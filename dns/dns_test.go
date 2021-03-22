package dns_test

import (
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logger := log.New(os.Stdout,"",log.Ltime)

	raftaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 7000,
	}

	httpaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 8000,
	}

	joinaddr := &net.TCPAddr{
		IP: net.ParseIP("1.2.3.4"),
		Port: 8000,
	}

	conf := cfg.Config{
		Hostname: "adin carik",
		LogLevel: 1,
		DataDir: "raft/test",
		Bootstrap: true,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: raftaddr,
		HttpAddr: httpaddr,
		JoinAddr: joinaddr,
	}
	node, err := raft.NewInMemNodeForTesting(&conf, logger)
	if err != nil{
		log.Fatal("Preparing tests failed: "+err.Error())
	}
	s := communication.HttpServer{
		Node: node,
		Address: httpaddr,
		Logger: logger,
	}
	go s.Start()

	//dns entrypoint
	dnsEntrypointLogger := *log.New(os.Stdout,"dns: ",log.Ltime)
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node: node,
		Address: conf.HttpAddr,
		Logger: &dnsEntrypointLogger,
	}
	dnsEntrypoint.Start()

	//dns api
	dnsApiLogger := *log.New(os.Stdout,"dns: ",log.Ltime)
	dnsApi := &communication.API{
		Node: node,
		//TODO: address and type from config
		Address: conf.HttpAddr,
		NetworkType: "udp",
		Logger: &dnsApiLogger,
	}
	go dnsApi.Start()
	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

type answerFormat struct {
	Domain string
	Error string
	Value string
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
		"Hostname": "dejan.ac.at",
		"Ip" : "1.2.3.4",
		"RequestType" : "store",
	}

	ans := SendBsonMessage("127.0.0.1:10000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Store failed: %s", ans)
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
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Update failed: %s", ans)
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
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		t.Errorf("Delete failed: %s", ans)
	}
}

/*func TestQuery(t *testing.T) {
	msg := bson.M{
		"Hostname": "www.ariel.dna",
		"Ip" : "2.4.8.10",
		"RequestType" : "store",
	}
	SendBsonMessage("127.0.0.1:10000",msg)

	ips, err := net.LookupIP("www.ariel.dna")
	if err != nil {
		t.Error(err)
	}
	if len(ips) < 1 {
		t.Errorf("No ip returned.")
	} else if !reflect.DeepEqual(ips[0], net.IP{2, 4, 8, 10}) {
		t.Errorf("Got wrong IP: %s", ips[0])
	}
}*/
