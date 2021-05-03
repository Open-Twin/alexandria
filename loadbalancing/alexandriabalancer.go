package loadbalancing

import (
	"context"
	"fmt"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"log"
	"net"
	"net/http"
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
		Nodes:     &lb.nodes,
		Interval:  lb.HealthCheckInterval,
		CheckType: HttpCheck,
	}
	hc.ScheduleHealthChecks()

	// https://gist.github.com/mike-zhang/3853251

	var ipport = 53
	var isport = 55
	var ishost = "localhost"
	//	var idrop *float64 = flag.Float64("d", 0.0, "Packet drop rate")

	hostport := fmt.Sprintf("%s:%d", ishost, isport)
	Vlogf(3, "Proxy port = %d, Server address = %s\n",
		ipport, hostport)

	if setup(hostport, ipport) {
		lb.runProxy()
	}
}

func (lb *AlexandriaBalancer) Close() {
	lb.httpServer.Shutdown(context.Background())
	// TODO Close udpServer and stop HealthChecks
}

func (lb *AlexandriaBalancer) startSignupEndpoint() {
	lb.httpServer = http.Server{Addr: ":8080"}
	http.HandleFunc("/signup", lb.addAlexandriaNode)
	err := lb.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Signup Endpoint for Loadbalancer could not be started: %s", err.Error())
	}
}

func (lb *AlexandriaBalancer) addAlexandriaNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) == 0 {
		log.Println("Reqeust without params")
		return
	}
	ip := r.Form["ip"][0]

	lb.lock.Lock()
	lb.nodes[ip] = dns.NodeHealth{
		Healthy:     false,
		Connections: 0,
	}
	lb.lock.Unlock()

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

func setup(hostport string, port int) bool {
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
	Vlogf(2, "Proxy serving on port %d\n", port)

	// Get server address
	srvaddr, err := net.ResolveUDPAddr("udp", hostport)
	if checkreport(1, err) {
		return false
	}
	ServerAddr = srvaddr
	Vlogf(2, "Connected to server at %s\n", hostport)
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
		Vlogf(3, "Relayed '%s' from server to %s.\n",
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
				Port: ServerAddr.Port,
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

	return &net.UDPAddr{
		Port: ServerAddr.Port,
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
		Vlogf(3, "Read '%s' from client %s\n",
			string(buffer[0:n]), cliaddr.String())
		saddr := cliaddr.String()
		dlock()
		conn, found := ClientDict[saddr]

		ServerAddr = lb.nextAddr()

		if !found {
			conn = NewConnection(ServerAddr, cliaddr)
			if conn == nil {
				dunlock()
				continue
			}
			ClientDict[saddr] = conn
			dunlock()
			Vlogf(2, "Created new connection for client %s\n", saddr)
			// Fire up routine to manage new connection
			go RunConnection(conn)
		} else {
			Vlogf(5, "Found connection for client %s\n", saddr)
			dunlock()
		}
		// Relay to server
		_, err = conn.ServerConn.Write(buffer[0:n])
		if checkreport(1, err) {
			continue
		}
	}
}

var verbosity = 6

// Log result if verbosity level high enough
func Vlogf(level int, format string, v ...interface{}) {
	if level <= verbosity {
		log.Printf(format, v...)
	}
}

// Handle errors
func checkreport(level int, err error) bool {
	if err == nil {
		return false
	}
	Vlogf(level, "Error: %s", err.Error())
	return true
}
