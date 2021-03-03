package raft

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"strconv"
	"strings"
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
		Hostname string `bson:"Hostname"`
		Ip string `bson:"Ip"`
		RequestType string `bson:"RequestType"`
	}{}

	if err := bson.Unmarshal(buf, &request); err != nil {
		logger.Println("Bad request: "+err.Error())
		return createResponse("","error","something went wrong. please check your input.")
	}

	//create new resource record
	rrecord := generateResourceRecord(request.Hostname, request.Ip)
	//marshal record
	event := dnsresource{
		Dnsormetadata: true,
		Hostname: request.Hostname,
		Ip: request.Ip,
		RequestType: request.RequestType,
		ResourceRecord: rrecord,
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
		resp = createResponse(request.Hostname,"ok","null")
	}
	return resp
}

func generateResourceRecord(hostname, ip string) dns.DNSResourceRecord {
	// split the string at the dots
	labels := strings.Split(hostname, ".")

	// reverse the slice order
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}

	//log.Println("data length: "+ strconv.Itoa(int(uint16(len([]byte(ip))))))
	//log.Println("ip: "+string([]byte(ip)))

	ipSections := strings.Split(ip, ".")
	ipByte := make([]byte, 4)
	for index, section := range ipSections {
		num, _ := strconv.Atoi(section)
		ipByte[index] = byte(num)
	}

	//TODO: values sind nur A records
	rrecord := dns.DNSResourceRecord{
		Labels:             labels,
		Type:               1, // A Record
		Class: 				1, // Internet Class
		TimeToLive:         100,
		ResourceDataLength: uint16(len(ipByte)),
		ResourceData:       ipByte,
	}
	return rrecord
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
