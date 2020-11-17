package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"net"
)

func main(){
	udpServer := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}
	udpServer2 := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 8333,
	}
	tcpServer := communication.TCPServer{
		Address: []byte{0,0,0,0},
		Port:    53,
	}
	quit := make(chan struct{})
	go udpServer.StartUDP(func (addr net.Addr, buf []byte) []byte{
		fmt.Println("MAIN: "+string(buf))
		return []byte("success")
	})
	go udpServer2.StartUDP(func (addr net.Addr, buf []byte) []byte{
		fmt.Println("MAIN: "+string(buf))
		return []byte("success2")
	})
	go tcpServer.StartTCP(func (addr net.Addr, buf []byte) []byte{
		fmt.Println("MAIN: "+string(buf))
		return []byte("success3")
	})
	<-quit
}