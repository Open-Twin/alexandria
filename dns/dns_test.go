package dns

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func SendBsonMessage(address string, msg bson.M) {
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
		fmt.Printf("Answer:\n%s\n", answer)
	} else {
		fmt.Printf("Error on receiving answer: %v", err)
	}
}

func TestStoreEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "store",
	}

	SendBsonMessage("127.0.0.1:10000",msg)
}

func TestUpdateEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "update",
	}

	SendBsonMessage("127.0.0.1:10000",msg)
}

func TestDeleteEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "delete",
	}

	SendBsonMessage("127.0.0.1:10000",msg)
}

func TestQuery(t *testing.T) {
	ips, err := net.LookupIP("www.dejan.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Query succesful: %v\n", ips)
}
