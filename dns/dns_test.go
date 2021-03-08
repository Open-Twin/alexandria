package dns

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
	reflreflect "reflect"
	"strings"
	"testing"
)

type answerFormat struct {
	Domain string
	Error string
	Value string
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func SendBsonMessage(address string, msg bson.M) string {
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

	return string(answer)
}

func TestStoreEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "store",
	}

	ans := SendBsonMessage("127.0.0.1:10000",msg)
	answerVals := answerFormat{}
	bson.Marshal(answerVals)
	if answerVals.Error != "ok" {
		t.Error("Store failed: %s", ans)
	}
}

func TestUpdateEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "update",
	}

	ans := SendBsonMessage("127.0.0.1:10000",msg)
	answerVals := answerFormat{}
	bson.Marshal(answerVals)
	if answerVals.Error != "ok" {
		t.Error("Store failed: %s", ans)
	}
}

func TestDeleteEntry(t *testing.T) {
	msg := bson.M{
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "delete",
	}

	ans := SendBsonMessage("127.0.0.1:10000",msg)
	answerVals := answerFormat{}
	bson.Marshal(answerVals)
	if answerVals.Error != "ok" {
		t.Error("Store failed: %s", ans)
	}
}

func TestQuery(t *testing.T) {
	ips, err := net.LookupIP("www.dejan.com")
	if err != nil {
		t.Error(err)
	}
	if !reflreflect.DeepEqual(ips[0], net.IP{1, 2, 3 ,4}) {
		t.Error("Got wrong IP: %s", ips[0])
	}
}
