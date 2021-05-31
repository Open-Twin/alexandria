package communication

import (
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/rs/zerolog/log"
	"net"
)

type DnsEntrypoint struct {
	Node    *raft.Node
	Address net.TCPAddr
}

func (api *DnsEntrypoint) Start(){
	udpserver := UDPServer{
		Address: api.Address.IP,
		Port: api.Address.Port,
	}
	tcpserver := TCPServer{
		Address: api.Address.IP,
		Port: api.Address.Port,
	}

	go udpserver.Start(func(addr net.Addr, buf []byte) []byte {
		answer := handle(addr,buf, api)
		return answer
	})
	go tcpserver.Start(func(addr net.Addr, buf []byte) []byte {
		answer := handle(addr,buf, api)
		return answer
	})
}

func handle(addr net.Addr, buf []byte, api *DnsEntrypoint) []byte{
	pdu := dns.HandleRequest(addr, buf)
	log.Debug().Msg("-------------------create answer-------------------")

	hostnames := dns.ExtractQuestionHostnames(&pdu)

	requestedRecords := queryDnsRepo(hostnames, api)
	log.Debug().Msgf("requested records: %v", requestedRecords)
	answer := dns.CreateAnswer(pdu, requestedRecords, buf)

	log.Debug().Msgf("Answer Header: %v\n", answer.Header)
	log.Debug().Msgf("Answer Flags: %v\n", answer.Flags)
	log.Debug().Msgf("Answer Answer Resource Records: %v\n", answer.AnswerResourceRecords)
	log.Debug().Msgf("Answer Additional Resource Records: %v\n", answer.AdditionalResourceRecords)
	log.Debug().Msgf("-------------------answer end-------------------")
	return dns.PrepareToSend(answer)
}

func queryDnsRepo(hostnames []string, api *DnsEntrypoint) []dns.DNSResourceRecord{
	array := make([]dns.DNSResourceRecord, 0)
	for _, hostname := range hostnames {
		query, err := api.Node.Fsm.DnsRepo.Read(hostname)
		if err != nil{
			//TODO recursion and error handling
			/*if pdu.Flags.RecursionDesired {
				//TODO: recursive lookup
				if pdu.
				recursiveAnswer, recErr := dns.RecursiveLookup(originalMessage)
				if recErr != nil {
					logging.Print(recErr.Error())
				}
			}*/
			log.Warn().Msgf("Requested domain not available: %s", hostname)
			return nil
		}else{
			array = append(array, query)
		}
	}
	return array
}
