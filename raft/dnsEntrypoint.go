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
		log.Println("-------------------create answer-------------------")
		answer := CreateAnswer(pdu, api.Node.fsm, api.Logger, buf)
		log.Println(answer.Header)
		log.Println(answer.Flags)
		log.Println(answer.AnswerResourceRecords)
		log.Println(string(answer.AnswerResourceRecords[0].ResourceData))
		log.Println("-------------------answer end-------------------")
		return PrepareToSend(answer)
	})
}
