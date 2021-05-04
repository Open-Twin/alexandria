package dns

import (
	"github.com/rs/zerolog/log"
)

func CreateAnswer(request DNSPDU, requestedRecords []DNSResourceRecord, originalMessage []byte) DNSPDU {
	// set Response Flag to true
	request.Flags.QueryResponse = true
	//Because server can handle recursion
	request.Flags.RecursionAvailable = true

	request.Flags.AuthoritativeAnswer = true

	request.Flags.AuthenticData = true

	request.Flags.CheckingDisabled = true

	if requestedRecords == nil {
		//TODO: other response codes
		request.Flags.ResponseCode = 3
		return request
	}

	answer := addAnswerResourceRecords(request, requestedRecords, originalMessage)
	return answer
}

func ExtractQuestionHostnames(pdu *DNSPDU) []string {
	hostnames := make([]string, 0)

	for _, question := range pdu.Questions {
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

func addAnswerResourceRecords(pdu DNSPDU, requestedRecords []DNSResourceRecord, originalMessage []byte) DNSPDU {
	for _, rrecord := range requestedRecords{
		pdu.Header.TotalAnswerResourceRecords += 1
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
		log.Error().Msgf("Error converting dns response to byte array: %v", err.Error())
	}
	return resp
}

func checkForPointer(originalLabels []string, records []DNSResourceRecord) []DNSResourceRecord{
	hostname := ConcatRevertLabels(originalLabels, false)
	for i, record := range records {
		if hostname == ConcatRevertLabels(record.Labels, true){
			//TODO: ??
			record.Labels = []string{"P", "O", "I", "N", "T", "E", "R"}
		}
		records[i] = record
	}
	return records
}
