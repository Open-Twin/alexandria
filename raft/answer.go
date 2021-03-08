package raft

import (
	"github.com/Open-Twin/alexandria/dns"
	"log"
)

var logging *log.Logger

func CreateAnswer(request dns.DNSPDU, fsm *Fsm, logger *log.Logger, originalMessage []byte) dns.DNSPDU {
	logging = logger

	// set Response Flag to true
	request.Flags.QueryResponse = true
	//Because server can handle recursion
	request.Flags.RecursionAvailable = true

	answer := addResourceRecords(request, fsm, originalMessage)

	return answer
}

func addResourceRecords(pdu dns.DNSPDU, fsm *Fsm, originalMessage []byte) dns.DNSPDU {
	for _, question := range pdu.Questions{
		pdu.Header.TotalAnswerResourceRecords += 1

		domainName := ""
		for i, part := range question.Labels {
			domainName += part
			if i < len(question.Labels)-1 {
				domainName += "."
			}
		}

		query, err := fsm.DnsRepo.Read(domainName)
		if err != nil{
			logging.Print(err.Error())
			/*if pdu.Flags.RecursionDesired {
				//TODO: recursive lookup
				if pdu.
				recursiveAnswer, recErr := dns.RecursiveLookup(originalMessage)
				if recErr != nil {
					logging.Print(recErr.Error())
				}
			}*/
		}
		log.Print("FABIAN:")
		log.Println(query)
		log.Print("--------")
		pdu.AnswerResourceRecords = append(pdu.AnswerResourceRecords, query)
	}
	return pdu
}

func PrepareToSend(pdu dns.DNSPDU) []byte {
	resp, err := pdu.Bytes()
	if err != nil{
		logging.Print(err.Error())
	}
	return resp
}
