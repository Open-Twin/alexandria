package main

import (
	"github.com/Open-Twin/alexandria/communication"
)

func main(){
	server := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}

	server.StartUDP()
}