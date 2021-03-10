package communication

import (
	"encoding/json"
	"errors"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/go-playground/validator/v10"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type API struct {
	Node    *raft.Node
	Address net.Addr
	NetworkType string
	Logger  *log.Logger
}

const ipv6Type = 28
const ipv4Type = 1

func (api *API) Start() {
	/*
	Metadata
	 */
	meta_udpserver := UDPServer{
		Address: []byte{0,0,0,0},
		Port: 20000,
	}

	log.Println("Starting DNS")
	go meta_udpserver.StartUDP(func(addr net.Addr, buf []byte) []byte {
		resp := handleMetadata(addr, buf, api.Node, api.Logger)
		return resp
	})

	/*
	DNS
	 */
	dns_udpserver := UDPServer{
		Address: []byte{0,0,0,0},
		Port: 10000,
	}

	log.Println("Starting DNS")
	go dns_udpserver.StartUDP(func(addr net.Addr, buf []byte) []byte {
		resp := handleDnsdata(addr, buf, api.Node, api.Logger)
		return resp
	})
}

func handleMetadata(addr net.Addr, buf []byte, node *raft.Node, logger *log.Logger) []byte{
	request := struct {
		Service string `bson:"Service"`
		Ip string `bson:"Ip"`
		Type string `bson:"Type"`
		Key string `bson:"Key"`
		Value string `bson:"Value"`
	}{}

	if err := bson.Unmarshal(buf, &request); err != nil {
		logger.Println("Bad request: "+err.Error())
		return createResponse("","error","something went wrong. please check your input.")
	}

	//marshal record
	event := storage.Metadata{
		Dnsormetadata: false,
		Service: request.Service,
		Ip: request.Ip,
		Type: request.Type,
		Key: request.Key,
		Value: request.Value,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		logger.Println("")
	}

	//Apply to Raft cluster
	applyFuture := node.RaftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		logger.Println("could not apply to raft cluster: "+err.Error())
		return createMetadataResponse(request.Service,request.Key,"error",err.Error())
	}
	var resp []byte
	if err != nil{
		resp = createMetadataResponse(request.Service,request.Key,"error","something went wrong. please check your input.")
	}else {
		resp = createMetadataResponse(request.Service,request.Key,"ok","null")
	}
	return resp
}

func createMetadataResponse(service, key, etype, value string) []byte{

	valueMap := map[string]string{
		"Type": etype,
		"Value": value,
	}
	response := struct {
		Service string
		Type string
		Key string
		Value map[string]string
	}{
		Service: service,
		Type: "response",
		Key: key,
		Value: valueMap,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Print("sendresponse failed")
	}

	return responseBytes
}

func handleDnsdata(addr net.Addr, buf []byte, node *raft.Node, logger *log.Logger) []byte{
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
	rrecord, err := generateResourceRecord(request.Hostname, request.Ip)
	//TODO: handle error

	//marshal record
	event := storage.Dnsresource{
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
	applyFuture := node.RaftNode.Apply(eventBytes, 5*time.Second)
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

func generateResourceRecord(hostname, ip string) (dns.DNSResourceRecord, error) {
	validate := validator.New()
	//check if hostname is valid
	errs := validate.Var(hostname,"required,hostname")
	//TODO: errors
	if errs != nil{
		return dns.DNSResourceRecord{}, errors.New("hostname is not valid")
	}
	//check if ip is valid ipv4 or ipv6
	errs = validate.Var(ip, "required,ipv4")
	ipType := ipv4Type
	if errs != nil{
		errs := validate.Var(ip, "required,ipv6")
		if errs != nil {
			return dns.DNSResourceRecord{}, errors.New("ip is not a valid ipv4 or ipv6 address")
		}
		ipType = ipv6Type
	}
	// split the string at the dots
	labels := strings.Split(hostname, ".")

	// reverse the slice order
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}


	ipSections := strings.Split(ip, ".")
	ipByte := make([]byte, 4)
	for index, section := range ipSections {
		num, _ := strconv.Atoi(section)
		ipByte[index] = byte(num)
	}

	//TODO:
	//RFC 1035: https://tools.ietf.org/html/rfc1035
	//RFC 2460 (ipv6): https://tools.ietf.org/html/rfc2460
	//Regex für ipv6:
	//(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))
	rrecord := dns.DNSResourceRecord{
		Labels:             labels,
		Type:               uint16(ipType), // A Record
		Class: 				1, // Internet Class
		TimeToLive:         1,
		ResourceDataLength: uint16(len(ipByte)),
		ResourceData:       ipByte,
	}
	return rrecord, nil
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