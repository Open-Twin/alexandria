package loadbalancing

import (
	"fmt"
	"net"
	//"github.com/Open-Twin/alexandria/communication"
	"../communication"
)

var dnsservers = [2]string{"192.168.0.160", "192.168.0.163"}
var pointer = 0
var dnsport = 8333

func loadbalanceAlexandriaNodes() {
	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    dnsport,
	}
	go udpServer.StartUDP(func(addr net.Addr, msg []byte) []byte {
		go forwardMsg(msg)
		return []byte("dns forwarded")
	})
}

func nextAddr() string {
	// Implementation of the loadbalancing
	pointer++
	if pointer > len(dnsservers) {
		pointer = 0
	}

	address := dnsservers[pointer]
	//adrentik := net.IPAddr{IP: net.ParseIP(address + string(dnsport))}
	//return &adrentik
	return address
}

func forwardMsg(msg []byte) {
	adrentik := nextAddr()
	dnsConn, err := net.Dial("udp", adrentik+":"+string(dnsport))
	if err != nil {
		fmt.Printf("We had an error.")
	}
	fmt.Fprintf(dnsConn, string(msg))
	dnsConn.Close()
	fmt.Print("Message forwareded to :", adrentik)
}
