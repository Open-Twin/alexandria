package dns

import (
	"github.com/Open-Twin/alexandria/communication"
	"log"
	"net"
)

func StartDNS(){
	server := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}

	log.Println("Starting DNS")
	server.StartUDP(func(addr net.Addr, buf []byte) []byte {
		ans := handleRequest(addr, buf)
		return ans
	})
}
