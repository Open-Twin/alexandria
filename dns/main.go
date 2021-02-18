package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Open-Twin/alexandria/communication"
	"log"
	"net"
)

func StartDNS(){
	server := communication.UDPServer{
		Address: []byte{0,0,0,0},
		Port: 53,
	}
	log.Println("Starting DNS")
	server.StartUDP(func(addr net.Addr, buf []byte) []byte {
		ans := handleRequest(addr, buf)
		return ans
	})
}

func handleRequest(addr net.Addr, buf []byte) []byte{
	/*Header*/
	header := buf[:13]
	id := header[:2]
	//flags
	flags := header[2:4]
	isResponse := flags[0] & 1
	//shift one bit to left to remove first bit
	opcode := flags[0] << 1
	//cancel out last 4 bits
	opcode = opcode >> 4
	opcode = opcode << 4
	authoritiveanswer := flags[0] & 32
	truncation := flags[0] & 64
	recursiondesired := flags[0] & 128
	recursionavailable := flags[1] & 1
	zfuture := flags[1] & 2
	authenticdata := flags[1] & 4
	checkingdisabled := flags[1] & 8
	rcode := flags[1] << 4
	//number of questions
	noq := header[4:6]
	//number of answers
	noa := header[6:8]
	//number of authorities
	noau := header[8:10]
	//number of additional records
	noar := header[10:12]

	/*Questions*/
	octat := buf[13]
	questions := buf[13:]
	cock := []byte("192.168.0.111")
	cock = append([]byte{byte(len(cock))},cock...)
	ans := append(append(buf[:12+octat],cock...),buf[12+octat:]...)
	ans[2] = ans[2] | 1
	//ans := buf[:12+octat] + cock + buf[12+octat+len(cock):]


	log.Println("-----HEADER-----")
	log.Println("-----LENGTH-----")
	log.Println(len(buf))
	log.Println("-----ID-----")
	log.Println(id)

	log.Println("-----FLAGS-----")
	log.Println("-----ISRESPONSE-----")
	log.Println(isResponse)
	log.Println("-----OPCODE-----")
	log.Println(opcode)
	log.Println("-----AUTHORITIVE-----")
	log.Println(authoritiveanswer)
	log.Println("-----TRUNCATION-----")
	log.Println(truncation)
	log.Println("-----RECDESIRED-----")
	log.Println(recursiondesired)
	log.Println("-----RECAVAIL-----")
	log.Println(recursionavailable)
	log.Println("-----ZFUTURE-----")
	log.Println(zfuture)
	log.Println("-----AUTHENTIC DATA-----")
	log.Println(authenticdata)
	log.Println("-----CHECKING DISABLED-----")
	log.Println(checkingdisabled)
	log.Println("-----RCODE-----")
	log.Println(rcode)

	log.Println("-----NUMBER OF QUESTIONS-----")
	log.Println(noq)
	log.Println("-----NUMBER OF ANSWERS-----")
	log.Println(noa)
	log.Println("-----NUMBER OF AUTHORITIES-----")
	log.Println(noau)
	log.Println("-----NUMBER OF ADDITIONAL RECORDS-----")
	log.Println(noar)

	log.Println("-----BODY-----")
	log.Println("-----OCTAT QUESTION 1-----")
	log.Println(octat)
	log.Println("-----Rest-----")
	log.Println(string(questions))
	return ans
}

func parseHeader(buffer *bytes.Buffer) (DNSHeader, error){
	var header DNSHeader
	err := binary.Read(buffer, binary.BigEndian, &header)

	if err != nil {
		return DNSHeader{}, errors.New("couldn't read DNS header")
	}

	return header, nil
}

func parseDnsMessage(buf []byte) (DNSPDU, error){
	buffer := bytes.NewBuffer(buf)

	header ,err := parseHeader(buffer)
	if err != nil{
		return DNSPDU{}, err
	}
	log.Println("--HEADER--")
	log.Println(header)

	return DNSPDU{},nil
}

