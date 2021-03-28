package main

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
)

func main() {
	if os.Args[1] == "dns" {
		sendDnsEntry()
	} else if os.Args[1] == "meta" {
		sendMetadataPacket()
	}
}

func sendBsonMessage(address string, msg bson.M) {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		log.Printf("Error on establishing connection: %s\n", err)
	}
	sendMsg, err := bson.Marshal(msg)

	conn.Write(sendMsg)
	log.Printf("Message sent: %s\n", sendMsg)

	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)
	if err != nil {
		log.Printf("Error on receiving answer: %v", err)
	} else {
		log.Printf("Answer:\n%s\n", answer)
	}
}

func sendDnsEntry() {
	address := "127.0.0.1:10000"
	/*msg := bson.M{
		"Labels":             []string{"at", "ac", "dejan"},
		"Type":               uint16(60),
		"Class":              uint16(50),
		"TimeToLive":         uint32(32),
		"ResourceDataLength": uint16(30),
		"ResourceData":       []byte{0x16, 0x32, 0x64},
		"RequestType":		  "store",
	}*/
	msg := bson.M{
		"Hostname":    os.Args[2],
		"Ip":          os.Args[3],
		"RequestType": "store",
	}

	sendBsonMessage(address, msg)
}

func sendMetadataPacket() {
	address := "127.0.0.1:20000"
	/*msg := bson.M{
		"Labels":             []string{"at", "ac", "dejan"},
		"Type":               uint16(60),
		"Class":              uint16(50),
		"TimeToLive":         uint32(32),
		"ResourceDataLength": uint16(30),
		"ResourceData":       []byte{0x16, 0x32, 0x64},
		"RequestType":		  "store",
	}*/
	msg := bson.M{
		"service": "toilet",
		"type":    os.Args[2],
		"key":     "temp",
		"value":   os.Args[3],
	}

	sendBsonMessage(address, msg)
}
