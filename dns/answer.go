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

func ExtractQuestionHostnames(pdu *DNSPDU) []string {
	hostnames := make([]string, 0)

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
	// find duplicate labels
	//if labels are duplicates, insert nil to mark a pointer
	pdu.AnswerResourceRecords =  checkForPointer(pdu.Questions[0].Labels,pdu.AnswerResourceRecords)
	pdu.AuthorityResourceRecords =  checkForPointer(pdu.Questions[0].Labels,pdu.AuthorityResourceRecords)
	pdu.AdditionalResourceRecords = checkForPointer(pdu.Questions[0].Labels,pdu.AdditionalResourceRecords)

	resp, err := pdu.Bytes()
	if err != nil{
		logging.Print(err.Error())
	}
	return resp
}
func checkForPointer(originalLabels []string, records []DNSResourceRecord) []DNSResourceRecord{
	hostname := ConcatRevertLabels(originalLabels, false)
	for i, record := range records {
		if hostname == ConcatRevertLabels(record.Labels, true){
			record.Labels = nil
		}
		records[i] = record
	}
	return records
}
