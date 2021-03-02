package raft

import (
	"github.com/Open-Twin/alexandria/dns"
	"log"
)

var logging *log.Logger

func CreateAnswer(request dns.DNSPDU, fsm *Fsm, logger *log.Logger) dns.DNSPDU {
	logging = logger

	// set Response Flag to true
	request.Flags.QueryResponse = true

	answer := addResourceRecords(request, fsm)

	return answer
}

func addResourceRecords(pdu dns.DNSPDU, fsm *Fsm) dns.DNSPDU {
	for _, question := range pdu.Questions{
		pdu.Header.TotalAnswerResourceRecords += 1

		domainName := ""
		for i, part := range question.Labels {
			domainName += part
			if i < len(question.Labels)-1 {
				domainName += "."
			}
		}

		log.Print("read dns")
		query, err := fsm.DnsRepo.Read(domainName)
		if err != nil{
			logging.Print(err.Error())
		}
		//log.Println("query dns: "+query.Labels[1])
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
