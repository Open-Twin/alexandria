package communication

import (
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/raft"
	"log"
	"net"
)

type DnsEntrypoint struct {
	Node    *raft.Node
	Address net.Addr
	Logger  *log.Logger
}

func (api *DnsEntrypoint) StartDnsEntrypoint(){
	udpserver := UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}
	tcpserver := TCPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}

	log.Println("Starting DNS entrypoint")
	go udpserver.StartUDP(func(addr net.Addr, buf []byte) []byte {
		answer := handle(addr,buf, api)
		return answer
	})
	go tcpserver.StartTCP(func(addr net.Addr, buf []byte) []byte {
		answer := handle(addr,buf, api)
		return answer
	})
}

func handle(addr net.Addr, buf []byte, api *DnsEntrypoint) []byte{
	pdu := dns.HandleRequest(addr, buf)
	log.Println("-------------------create answer-------------------")

	hostnames := dns.ExtractQuestionHostnames(pdu)
	log.Printf("HORST: %s", hostnames)

	requestedRecords := queryDnsRepo(hostnames, api)
	log.Printf("Ranshid: %s", requestedRecords)

	answer := dns.CreateAnswer(pdu, requestedRecords, api.Logger, buf)
	log.Println(answer.Header)
	log.Println(answer.Flags)
	log.Println(answer.AnswerResourceRecords)
	log.Println(string(answer.AnswerResourceRecords[0].ResourceData))
	log.Println("-------------------answer end-------------------")
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
			log.Printf("Reqeusted domain not available: %s", hostname)
		}else{
			array = append(array, query)
		}
	}
	return array
}
