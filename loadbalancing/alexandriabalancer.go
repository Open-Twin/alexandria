package loadbalancing

import (
	"context"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strings"
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

	// https://gist.github.com/mike-zhang/3853251

	//	var idrop *float64 = flag.Float64("d", 0.0, "Packet drop rate")

	dnsProxy := UdpProxy{
		Lb:   lb,
		Port: 53,
	}
	go dnsProxy.RunProxy()
	dnsApi := UdpProxy{
		Lb:   lb,
		Port: 10000,
	}
	go dnsApi.RunProxy()
	metaApi := UdpProxy{
		Lb:   lb,
		Port: 20000,
	}
	metaApi.RunProxy()
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
	ip := r.RemoteAddr

	// remove port from ip
	ip = ip[0:strings.LastIndex(ip, ":")]

	lb.lock.Lock()
	lb.nodes[ip] = dns.NodeHealth{
		Healthy:     true,
		Connections: 0,
	}
	lb.lock.Unlock()

	log.Info().Msgf("Node %s added", ip)
	w.Write([]byte("succesfully added"))
}

/**
Returns the next address in the list of loadbalanced nodes
*/
func (lb *AlexandriaBalancer) nextAddr() *net.UDPAddr {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	i := 0
	unhealthy := false
	for ip, health := range lb.nodes {
		if i == lb.pointer {
			if health.Healthy == true {
				return &net.UDPAddr{
					Port: lb.DnsPort,
					IP:   net.ParseIP(ip),
				}
			} else {
				unhealthy = true
			}
		}
		if unhealthy {
			if health.Healthy == true {
				return &net.UDPAddr{
					Port: lb.DnsPort,
					IP:   net.ParseIP(ip),
				}
			}
			if i == lb.pointer {
				break
			}
		}
		i += 1
		if i > len(lb.nodes) {
			i = 0
		}
	}

	lb.pointer = i

	return &net.UDPAddr{
		Port: lb.DnsPort,
		IP:   net.ParseIP("127.0.0.1"),
	}
}
