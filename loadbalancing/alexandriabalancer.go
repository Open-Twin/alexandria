package loadbalancing

import (
	"context"
	"fmt"
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

	var lbport = 53
	//	var idrop *float64 = flag.Float64("d", 0.0, "Packet drop rate")

	log.Info().Msgf("Loadbalancer port: %d\n", lbport)

	if setup(lbport) {
		lb.runProxy()
	}
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

// Information maintained for each client/server connection
type Connection struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}

// Generate a new connection by opening a UDP connection to the server
func NewConnection(srvAddr, cliAddr *net.UDPAddr) *Connection {
	conn := new(Connection)
	conn.ClientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if checkreport(1, err) {
		return nil
	}
	conn.ServerConn = srvudp
	return conn
}

// Global state
// Connection used by clients as the proxy server
var ProxyConn *net.UDPConn

// Address of server
var ServerAddr *net.UDPAddr

// Mapping from client addresses (as host:port) to connection
var ClientDict = make(map[string]*Connection)

// Mutex used to serialize access to the dictionary
var dmutex = new(sync.Mutex)

func setup(port int) bool {
	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if checkreport(1, err) {
		return false
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if checkreport(1, err) {
		return false
	}
	ProxyConn = pudp
	log.Info().Msgf("Loadbalancer listening on port %d\n", port)

	return true
}

func dlock() {
	dmutex.Lock()
}

func dunlock() {
	dmutex.Unlock()
}

// Go routine which manages connection from server to single client
func RunConnection(conn *Connection) {
	var buffer [1500]byte
	for {
		// Read from server
		n, err := conn.ServerConn.Read(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		// Relay it to client
		_, err = ProxyConn.WriteToUDP(buffer[0:n], conn.ClientAddr)
		if checkreport(1, err) {
			continue
		}
		log.Info().Msgf("Relayed '%s' from server to %s.\n",
			string(buffer[0:n]), conn.ClientAddr.String())
	}
}

/**
Returns the next address in the list of loadbalanced nodes
*/
func (lb *AlexandriaBalancer) nextAddr() *net.UDPAddr {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	i := 0
	for ip, health := range lb.nodes {
		if i == lb.pointer && health.Healthy == true {
			return &net.UDPAddr{
				Port: lb.DnsPort,
				IP:   net.ParseIP(ip),
			}
		} else {
			i += 1
		}
		if i > len(lb.nodes) {
			i = 0
			break
		}
	}

	lb.pointer = i

	return &net.UDPAddr{
		Port: lb.DnsPort,
		IP:   net.ParseIP("127.0.0.1"),
	}
}

// Routine to handle inputs to Proxy port
func (lb *AlexandriaBalancer) runProxy() {
	var buffer [1500]byte
	for {
		n, cliaddr, err := ProxyConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		log.Info().Msgf("Read '%s' from client %s\n", string(buffer[0:n]), cliaddr.String())
		saddr := cliaddr.String()
		dlock()
		conn, found := ClientDict[saddr]

		ServerAddr = lb.nextAddr()
		log.Info().Msgf("Forwarding node to %s", ServerAddr.String())

		if !found {
			conn = NewConnection(ServerAddr, cliaddr)
			if conn == nil {
				dunlock()
				continue
			}
			ClientDict[saddr] = conn
			dunlock()
			log.Info().Msgf("Created new connection for client %s\n", saddr)
			// Fire up routine to manage new connection
			go RunConnection(conn)
		} else {
			log.Info().Msgf("Found connection for client %s\n", saddr)
			dunlock()
		}
		// Relay to server
		_, err = conn.ServerConn.Write(buffer[0:n])
		if checkreport(1, err) {
			continue
		}
	}
}

// Handle errors
func checkreport(level int, err error) bool {
	if err == nil {
		return false
	}
	log.Error().Msgf("Error: %s", err.Error())
	return true
}
