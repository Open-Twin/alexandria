package raft

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net"
	"time"
)

type DnsApi struct {
	Node    *node
	Address net.Addr
	NetworkType string
	Logger  *log.Logger
}

func (api *DnsApi) StartDnsApi() {
	//ip := api.Address

	server := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 10000,
	}

	log.Println("Starting DNS")
	server.StartUDP(func(addr net.Addr, buf []byte) []byte {
		resp := handleData(addr, buf, api.Node, api.Logger)
		return resp
	})
}

func handleData(addr net.Addr, buf []byte, node *node, logger *log.Logger) []byte{

	request := struct {
		dnsormetadata bool
		record dns.DNSResourceRecord
	}{}

	if err := bson.Unmarshal(buf, &request); err != nil {
		logger.Println("Bad request")
		return createResponse("","error","something went wrong. please check your input.")
	}

	event := dnsresource{
		dnsormetadata: true,
		record: request.record,
	}
	log.Print("DEJANDNS: ")

	//query if domain exists
	domainName := ""
	for _, part := range request.record.Labels {
		domainName = part + domainName
	}
	exists := node.fsm.DnsRepo.Exists(domainName)
	if exists{
		return createResponse("","error","domain name already exists.")
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		logger.Println("")
	}
	//Apply to Raft cluster
	applyFuture := node.raftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		logger.Println("could not apply to raft cluster: "+err.Error())
		return createResponse("","error","something went wrong. please check your input.")
	}
	var resp []byte
	if err != nil{
		resp = createResponse("","error","something went wrong. please check your input.")
	}else {
		resp = createResponse(domainName,"ok","null")
	}
	return resp
}

func createResponse(domain, etype, value string) []byte{

	response := struct {
		Domain string
		Error string
		Value string
	}{
		Domain: domain,
		Error: etype,
		Value: value,
	}

	responseBytes, err := bson.Marshal(response)
	if err != nil {
		log.Print("sendresponse failed")
	}

	return responseBytes
}

