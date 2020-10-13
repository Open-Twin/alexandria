package main

import (
	"github.com/Open-Twin/alexandria/communication"
)

func main(){
	udpServer := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}
	tcpServer := communication.TCPServer{
		Address: []byte{0,0,0,0},
		Port:    53,
	}
	quit := make(chan struct{})
	go udpServer.StartUDP()
	go tcpServer.StartTCP()
	<-quit
}