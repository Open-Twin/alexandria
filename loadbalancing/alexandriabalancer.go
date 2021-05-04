package loadbalancing

import (
	"context"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type AlexandriaBalancer struct {
	DnsPort             int
	HealthCheckInterval int
	nodes               map[storage.Ip]dns.NodeHealth
	pointer             int
	lock                sync.RWMutex
	httpServer          http.Server
}

func (lb *AlexandriaBalancer) StartAlexandriaLoadbalancer() {
	lb.nodes = make(map[storage.Ip]dns.NodeHealth)
	lb.pointer = 0
	lb.lock = sync.RWMutex{}

	go lb.startSignupEndpoint()

	hc := HealthCheck{
		Nodes:     lb.nodes,
		Interval:  lb.HealthCheckInterval,
		CheckType: HttpCheck,
	}
	hc.ScheduleHealthChecks()

	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    lb.DnsPort,
	}

	// Listen for connections
	go udpServer.Start(func(addr net.Addr, msg []byte) []byte {
		// Run the method for every message received
		ans := lb.forwardMsg(addr, msg)
		return []byte(ans)
	})
}

func (lb *AlexandriaBalancer) Close() {
	log.Info().Msg("Shutting down Loadbalancer")
	lb.httpServer.Shutdown(context.Background())
	// TODO Close udpServer and stop HealthChecks
}

func (lb *AlexandriaBalancer) startSignupEndpoint() {
	lb.httpServer = http.Server{Addr: ":8080"}
	http.HandleFunc("/signup", lb.addAlexandriaNode)
	err := lb.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal().Msgf("Signup Endpoint for Loadbalancer could not be started: %s", err.Error())
	}
}

func (lb *AlexandriaBalancer) addAlexandriaNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ip := r.Form["ip"][0]

	lb.lock.Lock()
	lb.nodes[ip] = dns.NodeHealth{
		Healthy:     false,
		Connections: 0,
	}
	lb.lock.Unlock()

	w.Write([]byte("succesfully added"))
}

/**
Returns the next address in the list of loadbalanced nodes
*/
func (lb *AlexandriaBalancer) nextAddr() string {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	i := 0
	for ip, health := range lb.nodes {
		if i == lb.pointer && health.Healthy == true {
			return ip
		} else {
			i += 1
		}
		if i > len(lb.nodes) {
			i = 0
			break
		}
	}

	lb.pointer = i

	return ""
}

/**
Forwards an incoming message to a dns node
*/
func (lb *AlexandriaBalancer) forwardMsg(source net.Addr, msg []byte) string {
	log.Debug().Msgf("Message received: %s\n", string(msg))

	if len(lb.nodes) == 0 {
		log.Warn().Msg("No dns nodes to forward to.")
		return "no nodes available"
	}

	adrentik := lb.nextAddr()
	if adrentik == "" {
		log.Warn().Msg("No healthy nodes available.")
		return "No healthy nodes"
	}

	receiverAddr, err := net.ResolveUDPAddr("udp", adrentik+":"+strconv.Itoa(lb.DnsPort))
	if err != nil {
		log.Error().Msgf("Error on resolving dns address : %s\n", err)
	}

	sourceAddr, err := net.ResolveUDPAddr("udp", source.String())
	if err != nil {
		log.Error().Msgf("Error on resolving client address : %s\n", err)
	}

	target, err := net.DialUDP("udp", sourceAddr, receiverAddr)
	if err != nil {
		log.Error().Msgf("Error on establishing dns connection: %s\n", err)
	}

	_, err = target.WriteToUDP(msg, receiverAddr)
	if err != nil {
		log.Error().Msgf("Error on sending message to dns: %s\n", err)
	}

	log.Info().Msgf("Message forwareded to: %s:%d\n", adrentik, lb.DnsPort)
	return "yeah"
}
