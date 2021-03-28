package communication

import (
	"bufio"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/raft"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	raftaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 7000,
	}
	/*httpaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 8000,
	}*/
	metaaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 20000,
	}
	dnsaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 10000,
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
		MetaApiAddr: metaaddr,
		DnsApiAddr: dnsaddr,
		JoinAddr: joinaddr,
	}
	node, err := raft.NewInMemNodeForTesting(&conf)
	if err != nil{
		log.Fatal("Preparing tests failed: "+err.Error())
	}

	//dns api
	dnsApiLogger := *log.New(os.Stdout,"dns: ",log.Ltime)
	dnsApi := &API{
		Node: node,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress: conf.DnsApiAddr,
		NetworkType: "udp",
		Logger: &dnsApiLogger,
	}
	dnsApi.Start()

	/*var s = &HttpServer{
		Node:    node,
		Address: httpaddr,
		Logger:  logger,
	}
	go s.Start()*/


	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

type answerFormat struct {
	Service string `bson:"Service"`
	Type string	`bson:"Type"`
	Key string `bson:"Key"`
	Value map[string]string `bson:"Value"`
}

func SendBsonMessage(address string, msg bson.M) []byte {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		log.Printf("Error on establishing connection: %s\n", err)
	}
	sendMsg, err := bson.Marshal(msg)

	conn.Write(sendMsg)

	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)
	if err != nil {
		log.Printf("Error on receiving answer: %v", err)
	}
	return answer
}

func TestStoreEntryShouldPass(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "store",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["Type"] != "ok" {
		t.Errorf("Store failed: %s", ans)
	}
}

func TestUpdateEntryShouldPass(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "update",
		"key" : "voltage",
		"value" : "5",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Value["Type"] != "ok" {
		t.Errorf("Update failed: %s", ans)
	}
}

func TestDeleteEntryShouldPass(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "delete",
		"key" : "voltage",
		"value" : "5",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Value["Type"] != "ok" {
		t.Errorf("Delete failed: %s", ans)
	}
}