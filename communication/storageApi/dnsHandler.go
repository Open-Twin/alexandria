package storageApi

import (
	"encoding/json"
	"errors"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/go-playground/validator/v10"
	raftlib "github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

func handleDnsData(buf []byte, node *raft.Node) []byte {

	//handle other requests
	if node.RaftNode.State() != raftlib.Leader {
		resp, err := forwardToLeader(buf, string(node.RaftNode.Leader()), node.Config.DnsApiAddr.Port)
		if err != nil{
			log.Error().Msg("forward to leader failed")
			return communication.CreateResponse("", "error", "something went wrong. please check your input.")
		}
		return resp
	}

	request := struct {
		Hostname    string `bson:"hostname"`
		Ip          string `bson:"ip"`
		RequestType string `bson:"requestType"`
	}{}
	if err := bson.Unmarshal(buf, &request); err != nil {
		log.Error().Msgf("Bad request: %v", err.Error())
		return communication.CreateResponse("", "error", "something went wrong. please check your input.")
	}

	//create new resource record
	rrecord, err := generateResourceRecord(request.Hostname, request.Ip)
	//TODO: handle error

	//marshal record
	event := storage.Dnsresource{
		Dnsormetadata:  true,
		Hostname:       request.Hostname,
		Ip:             request.Ip,
		RequestType:    request.RequestType,
		ResourceRecord: rrecord,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error().Msg("unexpected failure")
	}
	//Apply to Raft cluster
	applyFuture := node.RaftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		log.Error().Msgf("could not apply to raft cluster: %v", err.Error())
		return communication.CreateResponse("", "error","something went wrong. please check your input")
	}
	if applyResp := applyFuture.Response(); applyResp != nil {
		switch v := applyResp.(type) {
		case error:
			log.Debug().Msgf("apply failed because of type %s",v.Error())
			log.Error().Msgf("could not apply to raft cluster: %v", applyResp.(error).Error())
			return communication.CreateResponse("", "error","something went wrong. please check your input")
		}
	}

	var resp []byte
	if err != nil {
		resp = communication.CreateResponse("", "error", "something went wrong. please check your input.")
	} else {
		resp = communication.CreateResponse(request.Hostname, "ok", "null")
	}
	return resp
}

func generateResourceRecord(hostname, ip string) (dns.DNSResourceRecord, error) {
	validate := validator.New()
	//check if hostname is valid
	errs := validate.Var(hostname, "required,hostname")
	//TODO: errors
	if errs != nil {
		return dns.DNSResourceRecord{}, errors.New("hostname is not valid")
	}
	//check if ip is valid ipv4 or ipv6
	errs = validate.Var(ip, "required,ipv4")
	ipType := ipv4Type
	if errs != nil {
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
		Class:              1,              // Internet Class
		TimeToLive:         1,
		ResourceDataLength: uint16(len(ipByte)),
		ResourceData:       ipByte,
	}
	return rrecord, nil
}