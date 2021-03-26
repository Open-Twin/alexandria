package dns_test

import (
	"context"
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const apiAddr = "127.0.0.1"
const apiPort = 10001
const entrypointAddr = "127.0.0.1"
const entrypointPort = 10002

func TestMain(m *testing.M) {
	logger := log.New(os.Stdout, "", log.Ltime)

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
	node, err := raft.NewInMemNodeForTesting(&conf, logger)
	if err != nil {
		log.Fatal("Preparing tests failed: " + err.Error())
	}
	/*s := communication.HttpServer{
		Node: node,
		Address: httpaddr,
		Logger: logger,
	}
	go s.Start()*/

	//dns entrypoint
	dnsEntrypointLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node:    node,
		Address: conf.DnsAddr,
		Logger:  &dnsEntrypointLogger,
	}
	dnsEntrypoint.Start()

	//dns api
	dnsApiLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	dnsApi := &communication.API{
		Node: node,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress:  conf.DnsApiAddr,
		NetworkType: "udp",
		Logger:      &dnsApiLogger,
	}
	go dnsApi.Start()
	time.Sleep(5 * time.Second)

	code := m.Run()
	os.Exit(code)
}

type answerFormat struct {
	Domain string
	Error  string
	Value  string
}

func SendBsonMessage(address string, msg bson.M) []byte {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		fmt.Printf("Error on establishing connection: %s\n", err)
	}
	sendMsg, _ := bson.Marshal(msg)

	conn.Write(sendMsg)
	fmt.Printf("Message sent: %s\n", sendMsg)

	answer := make([]byte, 2048)
	_, err = conn.Read(answer)

	if err != nil {
		fmt.Printf("Error on receiving answer: %v", err.Error())
	} else {
		fmt.Printf("Answer:\n%s\n", answer)
	}

	return answer
}

func TestStoreEntry(t *testing.T) {
	msg := bson.M{
		"Hostname":    "dejan.ac.at",
		"Ip":          "1.2.3.4",
		"RequestType": "store",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Store failed: %s", ans)
	}
}

func TestUpdateEntry(t *testing.T) {
	msg := bson.M{
		"Hostname":    "dejan.ac.at",
		"Ip":          "1.2.3.4",
		"RequestType": "update",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Update failed: %s", ans)
	}
}

func TestDeleteEntry(t *testing.T) {
	msg := bson.M{
		"Hostname":    "dejan.ac.at",
		"Ip":          "1.2.3.4",
		"RequestType": "delete",
	}

	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Delete failed: %s", ans)
	}
}

func TestQuery(t *testing.T) {
	msg := bson.M{
		"Hostname":    "www.ariel.dna",
		"Ip":          "2.4.8.10",
		"RequestType": "store",
	}
	ans := SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	log.Printf("Storing ariel: %s", ans)

	ips, err := sendDig("www.ariel.dna", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if err != nil {
		t.Error(err)
	}
	if len(ips) < 1 {
		t.Errorf("No ip returned.")
	} else if !reflect.DeepEqual(ips[0], net.IP{2, 4, 8, 10}) {
		t.Errorf("Got wrong IP: %s", ips[0])
	}
}

func TestDnsNodeDistribution(t *testing.T) {
	msg := bson.M{
		"Hostname":    "eenie.meenie",
		"Ip":          "1.1.1.1",
		"RequestType": "store",
	}
	SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	ip, _ := sendDig("eenie.meenie", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if len(ip) < 1 {
		t.Errorf("No ip returned.")
	}
	if !reflect.DeepEqual(ip[0], net.IP{1, 1, 1, 1}) {
		t.Errorf("Loadbalancer returned %s instead of 1.1.1.2", ip)
	}

	msg = bson.M{
		"Hostname":    "eenie.meenie",
		"Ip":          "1.1.1.2",
		"RequestType": "store",
	}
	SendBsonMessage(apiAddr+":"+strconv.Itoa(apiPort), msg)
	ip, _ = sendDig("eenie.meenie", entrypointAddr+":"+strconv.Itoa(entrypointPort))
	if reflect.DeepEqual(ip[0], net.IP{1, 1, 1, 2}) {
		t.Errorf("Loadbalancer returned %s instead of 1.1.1.2", ip)
	}
}

func sendDig(hostname string, dnsip string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, dnsip)
		},
	}
	return r.LookupIP(context.Background(), "ip", hostname)
}
