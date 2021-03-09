package dns

import (
	"log"
)

var logging *log.Logger

func CreateAnswer(request DNSPDU, requestedRecords []DNSResourceRecord, logger *log.Logger, originalMessage []byte) DNSPDU {
	logging = logger

	// set Response Flag to true
	request.Flags.QueryResponse = true
	//Because server can handle recursion
	request.Flags.RecursionAvailable = true

	request.Flags.AuthoritativeAnswer = true

	request.Flags.AuthenticData = true

	request.Flags.CheckingDisabled = true

	answer := addResourceRecords(request, requestedRecords, originalMessage)

	return answer
}

func ExtractQuestionHostnames(pdu DNSPDU) []string {
	hostnames := []string{}

	for _, question := range pdu.Questions {
		pdu.Header.TotalAnswerResourceRecords += 1

		domainName := ""
		for i, part := range question.Labels {
			domainName += part
			if i < len(question.Labels)-1 {
				domainName += "."
			}
		}

		hostnames = append(hostnames, domainName)
	}

	return hostnames
}

func addResourceRecords(pdu DNSPDU, requestedRecords []DNSResourceRecord, originalMessage []byte) DNSPDU {
	for _, rrecord := range requestedRecords{
		pdu.AnswerResourceRecords = append(pdu.AnswerResourceRecords, rrecord)
	}
	return pdu
}

func PrepareToSend(pdu DNSPDU) []byte {
	resp, err := pdu.Bytes()
	if err != nil{
		logging.Print(err.Error())
	}
	return resp
}
