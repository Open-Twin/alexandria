package communication

import (
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/raft"
	"log"
	"net"
)

type DnsEntrypoint struct {
	Node    *raft.Node
	Address net.TCPAddr
	Logger  *log.Logger
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

	log.Println("Starting DNS entrypoint")
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
	log.Println("-------------------create answer-------------------")

	hostnames := dns.ExtractQuestionHostnames(pdu)
	log.Printf("HORST: %s", hostnames)

	requestedRecords := queryDnsRepo(hostnames, api)

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
