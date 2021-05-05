package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
)

func main() {
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
		"hostname":    os.Args[1],
		"ip":          os.Args[2],
		"requestType": "store",
	}

	sendBsonMessage(address, msg)
}

func sendBsonMessage(address string, msg bson.M) {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		fmt.Printf("Error on establishing connection: %s\n", err)
	}
	sendMsg, err := bson.Marshal(msg)

	conn.Write(sendMsg)
	fmt.Printf("Message sent: %s\n", sendMsg)

	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)
	if err != nil {
		fmt.Printf("Error on receiving answer: %v", err)
	} else {
		fmt.Printf("Answer:\n%s\n", answer)
	}
}
