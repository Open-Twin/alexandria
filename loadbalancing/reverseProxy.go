package loadbalancing

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"sync"
)

// Information maintained for each client/server connection
type Connection struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}

type UdpProxy struct {
	// Connection used by clients as the proxy server
	proxyConn *net.UDPConn
	Lb        *AlexandriaBalancer
	Port      int
	// Address of server
	serverAddr *net.UDPAddr
	// Mapping from client addresses (as host:port) to connection
	clientDict map[string]*Connection
	// Mutex used to serialize access to the dictionary
	dmutex sync.Mutex
}

// Routine to handle inputs to Proxy port
func (up *UdpProxy) RunProxy() {
	if !up.setup(up.Port) {
		return
	}

	var buffer [1500]byte
	for {
		n, cliaddr, err := up.proxyConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		log.Info().Msgf("Read '%s' from client %s\n", string(buffer[0:n]), cliaddr.String())
		saddr := cliaddr.String()
		up.dlock()
		conn, found := up.clientDict[saddr]

		up.serverAddr = up.Lb.nextAddr(up.Port)
		log.Info().Msgf("Forwarding node to %s", up.serverAddr.String())

		if !found {
			conn = newConnection(up.serverAddr, cliaddr)
			if conn == nil {
				up.dunlock()
				continue
			}
			up.clientDict[saddr] = conn
			up.dunlock()
			log.Info().Msgf("Created new connection for client %s\n", saddr)
			// Fire up routine to manage new connection
			go up.runConnection(conn)
		} else {
			log.Info().Msgf("Found connection for client %s\n", saddr)
			up.dunlock()
		}
		// Relay to server
		log.Info().Msg("NIGNOG")
		_, err = conn.ServerConn.Write(buffer[0:n])
		if checkreport(1, err) {
			continue
		}
	}
}

// Generate a new connection by opening a UDP connection to the server
func newConnection(srvAddr, cliAddr *net.UDPAddr) *Connection {
	conn := new(Connection)
	conn.ClientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if checkreport(1, err) {
		return nil
	}
	conn.ServerConn = srvudp
	return conn
}

func (up *UdpProxy) setup(port int) bool {
	up.clientDict = make(map[string]*Connection)

	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if checkreport(1, err) {
		return false
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if checkreport(1, err) {
		return false
	}
	up.proxyConn = pudp
	log.Info().Msgf("Loadbalancer listening on port %d\n", port)

	return true
}

func (up *UdpProxy) dlock() {
	up.dmutex.Lock()
}

func (up *UdpProxy) dunlock() {
	up.dmutex.Unlock()
}

// Go routine which manages connection from server to single client
func (up *UdpProxy) runConnection(conn *Connection) {
	var buffer [1500]byte
	for {
		// Read from server
		n, err := conn.ServerConn.Read(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		log.Info().Msg("NIGNOG2")
		// Relay it to client
		_, err = up.proxyConn.WriteToUDP(buffer[0:n], conn.ClientAddr)
		if checkreport(1, err) {
			continue
		}
		log.Info().Msgf("Relayed '%s' from server to %s.\n",
			string(buffer[0:n]), conn.ClientAddr.String())
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
