package storageApi

import (
	"bufio"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"strings"
)

type API struct {
	Node        *raft.Node
	MetaAddress net.TCPAddr
	DNSAddress  net.TCPAddr
	NetworkType string
}

const ipv6Type = 28
const ipv4Type = 1

func (api *API) Start() {
	/*
		Metadata
	*/
	meta_udpserver := communication.UDPServer{
		Address: api.MetaAddress.IP,
		Port:    api.MetaAddress.Port,
	}

	log.Info().Msg("Starting DNS")
	go meta_udpserver.Start(func(addr net.Addr, buf []byte) []byte {
		log.Info().Msg("got metadata request")
		resp := handleMetadata(buf, api.Node)
		return resp
	})

	/*
		DNS
	*/
	dns_udpserver := communication.UDPServer{
		Address: api.DNSAddress.IP,
		Port:    api.DNSAddress.Port,
	}

	log.Info().Msg("Starting DNS API")
	go dns_udpserver.Start(func(addr net.Addr, buf []byte) []byte {

		resp := handleDnsData(buf, api.Node)
		return resp
	})
}

func forwardToLeader(eventBytes []byte, leaderAddr string, port int) ([]byte, error) {
	addr := strings.Split(leaderAddr, ":")[0]
	leader := addr + ":" + strconv.Itoa(port)
	log.Info().Msg("forwarding request to leader: " + leader)
	con, err := net.Dial("udp", leader)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	//write to leader
	_, err = con.Write(eventBytes)
	if err != nil {
		return nil, err
	}
	//read answer
	p := make([]byte, 2048)
	_, err = bufio.NewReader(con).Read(p)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("Request forwarded to leader %s", leader)

	return p, nil
}
