package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

var request []byte

func HandleRequest(addr net.Addr, requestData []byte) DNSPDU {
	request = requestData

	// Create buffer from request array
	buf := bytes.NewBuffer(request)

	// Parse header of dns message
	header, err := parseHeader(buf)
	if err != nil {
		log.Print(err.Error())
	}

	// Parse flags of dns header
	flags, err := parseFlags(header.Flags)
	if err != nil {
		log.Print(err.Error())
	}

	// Parse dns body containing Question, Answer, Authority and Additional Information
	body, err := parseBody(header, buf)
	if err != nil {
		log.Print(err.Error())
	}

	body.Flags = flags

	log.Printf("------Header------\n %+v \n", header)
	log.Printf("------Flags-------\n %+v \n", flags)
	log.Printf("------Body--------\n %+v \n", body)


	// Send answer
	return body
}

/*
Reads the dns header from the byte buffer into the struct
 */
func parseHeader(buffer *bytes.Buffer) (DNSHeader, error){
	var header DNSHeader

	// Read the bytes into the struct
	err := binary.Read(buffer, binary.BigEndian, &header) // BigEndian is the order of transmission used by the dns protocol
	if err != nil {
		return DNSHeader{}, errors.New("Error reading DNS header: "+err.Error())
	}

	return header, nil
}

/*
Reads the dns flags out of the 2 byte long flag section in the header.
 */
func parseFlags(header uint16) (DNSFlags, error) {
	var flags DNSFlags
	// 1 bit long query response flag (1st bit)
	// The value is shifted 15 bits to the right to remove all other bits
	if header >> 15 == 1 {
		flags.QueryResponse = true
	}

	// 4 bit long OPCODE
	// The value is shifted one bit to the left to remove all the left bits. It is then shifted 12 bits to the right to remove all the rights bits.
	flags.OpCode = uint8(header << 1 >> 12)

	// 1 bit long authorative answer (6th bit)
	// The value is shifted five bits to the left to remove all the left bits. It is then shifted 15 bits to the right to remove all the rights bits.
	if header << 5 >> 15 == 1 {
		flags.AuthoritativeAnswer = true
	}

	// 1 bit long truncation (7th bit)
	if header << 6 >> 15 == 1 {
		flags.Truncated = true
	}

	// 1 bit long recursion desired (8th bit)
	if header << 7 >> 15 == 1 {
		flags.RecursionDesired = true
	}

	// 1 bit long recursion available  (9th bit)
	if header << 8 >> 15 == 1 {
		flags.RecursionAvailable = true
	}

	//z
	if header << 9 >> 15 == 1 {
		flags.Z = true
	}

	//ad
	if header << 10 >> 15 == 1 {
		flags.AuthenticData = true
	}

	//cd
	if header << 11 >> 15 == 1 {
		flags.CheckingDisabled = true
	}

	// 4 bit long response code
	flags.ResponseCode = uint8(header << 12 >> 12)

	return flags, nil
}

/*
Reads questions, answers, authorities and additional information from the query body
https://docstore.mik.ua/orelly/networking_2ndEd/dns/appa_02.htm
*/
func parseBody(header DNSHeader, buffer *bytes.Buffer) (DNSPDU, error) {
	pdu := DNSPDU{}
	//add header to pdu
	pdu.Header = header

	/* read question section */
	questions := make([]DNSQuestion, 0)
	//loops over question count and reads the questions
	for i := 0; i < int(pdu.Header.TotalQuestions); i++ {
		question := DNSQuestion{}

		labels, err := readLabels(buffer)
		question.Labels = labels

		if err != nil {
			return DNSPDU{}, err
		}

		question.Type = binary.BigEndian.Uint16(buffer.Next(2))
		question.Class = binary.BigEndian.Uint16(buffer.Next(2))

		questions = append(questions, question)
	}
	//add question section to pdu
	pdu.Questions = questions

	/* read answer, authority and additional section format */
	answer, err := readResourceRecords(buffer, int(header.TotalAnswerResourceRecords))
	authority, err := readResourceRecords(buffer, int(header.TotalAuthorityResourceRecords))
	additional, err := readResourceRecords(buffer, int(header.TotalAdditionalResourceRecords))
	if err != nil {
		return DNSPDU{}, err
	}
	pdu.AnswerResourceRecords = answer
	pdu.AuthorityResourceRecords = authority
	pdu.AdditionalResourceRecords = additional

	return pdu, nil
}
/*
 * reads the resource records of a DNS message and returns them
 */
func readResourceRecords(buffer *bytes.Buffer, resourceCount int) ([]DNSResourceRecord, error) {
	var resourceRecords []DNSResourceRecord
	for i := 0; i < resourceCount; i++ {
		resourceRecord := DNSResourceRecord{}

		//Name
		labels, err := readLabels(buffer)
		if err != nil {
			return []DNSResourceRecord{}, nil
		}
		resourceRecord.Labels = labels

		//type
		resourceRecord.Type = binary.BigEndian.Uint16(buffer.Next(2))

		//class
		resourceRecord.Class = binary.BigEndian.Uint16(buffer.Next(2))

		//Time to live
		resourceRecord.TimeToLive = binary.BigEndian.Uint32(buffer.Next(4))

		//Resource Data Length
		resourceRecord.ResourceDataLength = binary.BigEndian.Uint16(buffer.Next(2))

		//Resource Data
		resourceRecord.ResourceData = buffer.Next(int(resourceRecord.ResourceDataLength))

		resourceRecords = append(resourceRecords, resourceRecord)
	}

	return resourceRecords, nil
}
/*
 * reads the labels of a qname (question or answer section) and returns them
 * it distinguishes between dns pointers and actual labels
 */
func readLabels(buffer *bytes.Buffer) ([]string, error){
	var labels []string
	//loops over the length of the labels
	//the first octet (byte) is the length which is followed by that number of octets
	//the domain name terminates with a zero length octet
	for length := buffer.Next(1); length[0] != 0; length = buffer.Next(1) {
		//distinguish between dns pointer or label
		//a pointer starts with two leading 1's
		//a label starts with two leading 0's
		if length[0] >> 6 == 3 {
			//POINTER
			//DNS Pointers point to a specific byte in the whole message
			//pointers consist of two bytes, which is why the second byte is important
			pointerSecondByte := buffer.Next(1)
			pointer := []byte{length[0] << 2 >> 2, pointerSecondByte[0]}
			//get the byte the pointer is pointing to
			labelStart := binary.BigEndian.Uint16(pointer)

			var labelBytes []byte
			//append the labels that the pointer is pointing to
			for i := labelStart; request[i] != 0; i++ {
				labelBytes = append(labelBytes, request[i])
			}
			labels = append(labels, string(labelBytes))

		}else if length[0] >> 6 == 0 {
			//LABEL
			//get length of labels
			labelLength := int(length[0])
			//get the label
			labelBytes := buffer.Next(labelLength)
			//append it to the labels
			labels = append(labels, string(labelBytes))
		}
	}

	/*for _, label := range labels {
		log.Printf("%s\n", label)
	}*/

	return labels, nil
}