package dns_test

import (
	"bufio"
	"context"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/communication/storageApi"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/rs/zerolog"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

const apiAddr = "127.0.0.1"
const apiPort = 10001
const entrypointAddr = "127.0.0.1"
const entrypointPort = 10002

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Level(3))
	raftaddr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 7001,
	}

	httpaddr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8001,
	}

	dnsaddr := net.TCPAddr{
		IP:   net.ParseIP(entrypointAddr),
		Port: entrypointPort,
	}

	dnsapiaddr := net.TCPAddr{

		IP:   net.ParseIP(apiAddr),
		Port: apiPort,
	}

	conf := cfg.Config{
		Hostname:            "adin carik",
		LogLevel:            1,
		DataDir:             "raft/test",
		Bootstrap:           true,
		Autojoin:            false,
		HealthcheckInterval: 2000,
		RaftAddr:            raftaddr,
		HttpAddr:            httpaddr,
		DnsApiAddr:          dnsapiaddr,
		DnsAddr:             dnsaddr,
	}
	node, err := raft.NewInMemNodeForTesting(&conf)
	if err != nil {
		log.Printf("Error Creating InMemoryNodeForTesting: %v", err.Error())
	}

	/*s := communication.HttpServer{
		Node: node,
		Address: httpaddr,
	}
	go s.Start()*/

	//dns entrypoint
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node:    node,
		Address: conf.DnsAddr,
	}
	dnsEntrypoint.Start()

	//dns api
	dnsApi := &storageApi.API{
		Node: node,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress:  conf.DnsApiAddr,
		NetworkType: "udp",
	}
	go dnsApi.Start()
	time.Sleep(5 * time.Second)

	code := m.Run()
	os.Exit(code)
}

type answerFormat struct {
	Domain string `bson:"domain"`
	Error  string `bson:"error"`
	Value  string `bson:"value"`
}

func SendBsonMessage(address string, msg bson.M, t *testing.T) []byte {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		t.Errorf("Error on establishing connection: %s\n", err)
	}
	sendMsg, _ := bson.Marshal(msg)

	conn.Write(sendMsg)
	t.Logf("Message sent: %s\n", sendMsg)

	answer := make([]byte, 2048)
	_, err = conn.Read(answer)

	if err != nil {
		t.Errorf("Error on receiving answer: %v", err.Error())
	} else {
		t.Logf("Answer:\n%s\n", answer)
	}

	return answer
}

/*
	TESTING THE DNS AVAILABILITY
*/

func TestDNSShouldBeReachable(t *testing.T) {
	conn, err := net.DialTimeout("udp", apiAddr+":"+strconv.Itoa(apiPort), 500)
	defer conn.Close()
	if err != nil {
		t.Errorf("Error on establishing connection: %s\n", err)
	}

	timeoutDuration := 5 * time.Second
	conn.SetDeadline(time.Now().Add(timeoutDuration))

	msg := []byte("hello")
	conn.Write(msg)
	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)

	if err != nil {
		t.Errorf("Test failed. Cannot read answer")
	}
}

func TestDNSShouldNotBeReachable(t *testing.T) {
	conn, err := net.DialTimeout("udp", apiAddr+":9999", 500)
	defer conn.Close()
	if err != nil {
		t.Errorf("Error on establishing connection: %s\n", err)
	}

	timeoutDuration := 5 * time.Second
	conn.SetDeadline(time.Now().Add(timeoutDuration))

	msg := []byte("hello")
	conn.Write(msg)
	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)

	if err == nil {
		t.Errorf("Test failed. Can read answer")
	}
}

/*
	TESTING THE DNS FUNCTIONALITY
*/

func TestStoreEntryPass(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "store",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Store failed: %s", ans)
	}
}

func TestStoreEntryFail(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "save",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "error" {
		t.Errorf("Store should not go through, due to wrong type: %s", ans)
	}
}

func TestUpdateEntryPass(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "update",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Update failed: %s", ans)
	}
}

func TestUpdateEntryFail(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "modify",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "error" {
		t.Errorf("Update should not go through, due to wrong type: %s", ans)
	}
}

func TestDeleteEntryPass(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "delete",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Delete failed: %s", ans)
	}
}

func TestDeleteEntryFail(t *testing.T) {
	msg := bson.M{
		"hostname":    "dejan.ac.at",
		"ip":          "1.2.3.4",
		"requestType": "remove",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "error" {
		t.Errorf("Delete should not go through, due to wrong type: %s", ans)
	}
}

func TestQuery(t *testing.T) {
	gesuchtip := "2.4.8.10"
	msg := bson.M{
		"hostname":    "www.ariel.dna",
		"ip":          gesuchtip,
		"requestType": "store",
	}
	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	t.Logf("Storing ariel: %s", ans)

	ips, err := sendDig("www.ariel.dna", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if err != nil {
		t.Error(err)
	}
	if len(ips) < 1 {
		t.Errorf("No ip returned.")
	} else if ips[0].String() != gesuchtip {
		t.Errorf("Got wrong IP: %s", ips[0])
	}
}

func TestDnsNodeDistribution(t *testing.T) {
	msg := bson.M{
		"hostname":    "eenie.meenie",
		"ip":          "1.1.1.1",
		"requestType": "store",
	}
	SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	time.Sleep(1 * time.Second)
	ip1, _ := sendDig("eenie.meenie", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if len(ip1) < 1 {
		t.Errorf("No ip returned.")
		return
	}

	msg = bson.M{
		"hostname":    "eenie.meenie",
		"ip":          "1.1.1.2",
		"requestType": "store",
	}
	SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg, t)
	time.Sleep(1 * time.Second)
	ip2, _ := sendDig("eenie.meenie", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if len(ip2) < 1 {
		t.Errorf("No ip returned.")
		return
	}

	if ip1[0].String() == ip2[0].String() {
		t.Errorf("Same ip was returned two times: %s, %s", ip1[0], ip2[0])
	}
}

func sendDig(hostname string, dnsip string) ([]net.IPAddr, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, dnsip)
		},
	}
	return r.LookupIPAddr(context.Background(), hostname)
}
