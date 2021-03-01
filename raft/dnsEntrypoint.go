package raft

import (
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"log"
	"net"
)

type DnsEntrypoint struct {
	Node    *node
	Address net.Addr
	Logger  *log.Logger
}

func (api *DnsEntrypoint) StartDnsEntrypoint(){
	server := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}

	log.Println("Starting DNS entrypoint")
	server.StartUDP(func(addr net.Addr, buf []byte) []byte {
		pdu := dns.HandleRequest(addr, buf)

		answer := CreateAnswer(pdu, api.Node.fsm, api.Logger)

		return PrepareToSend(answer)
	})
}
