package storageApi

import (
	"bufio"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/raft"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

const leaderaddr = "leader"
const follower1 = "follower1"

func TestMain(m *testing.M) {

	raftaddr := net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 7000,
	}

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

	dnsApi := &API{
		Node: node,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress: conf.DnsApiAddr,
		NetworkType: "udp",
	}
	dnsApi.Start()

	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

type answerFormat struct {
	Service string `bson:"service"`
	Type string	`bson:"type"`
	Key string `bson:"key"`
	Value map[string]string `bson:"value"`
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

/*
	TESTING THE APIs AVAILABILITY
*/

func TestAPIShouldBeReachable(t *testing.T) {
	conn, err := net.DialTimeout("udp", "127.0.0.1:20000",500)
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

func TestAPIShouldNotBeReachable(t *testing.T) {
	conn, err := net.DialTimeout("udp", "127.0.0.1:20001",500)
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
	TESTING THE APIs FUNCTIONALITY
 */

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
	if answerVals.Value["type"] != "ok" {
		t.Errorf("Store failed: %s", ans)
	}
}

func TestStoreEntryShouldFail(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "save",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["type"] != "error" {
		t.Errorf("Store should not go through, due to wrong type: %s", ans)
	}
}

func TestGetEntryShouldPass(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "get",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["type"] != "data" {
		t.Errorf("Get failed: %s", ans)
	}
}

func TestGetEntryShouldFail(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "receive",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["type"] != "error" {
		t.Errorf("Get shouldnt go throughm, due to wrong type: %s", ans)
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
	if answerVals.Value["type"] != "ok" {
		t.Errorf("Update failed: %s", ans)
	}
}

func TestUpdateEntryShouldFail(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "modify",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["type"] != "error" {
		t.Errorf("Update should not go through, due to wrong type: %s", ans)
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
	if answerVals.Value["type"] != "ok" {
		t.Errorf("Delete failed: %s", ans)
	}
}

func TestDeleteEntryShouldFail(t *testing.T) {
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "remove",
		"key" : "voltage",
		"value" : "3",
	}

	ans := SendBsonMessage("127.0.0.1:20000",msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Print(answerVals)
	if answerVals.Value["type"] != "error" {
		t.Errorf("Delete should not go through, due to wrong type: %s", ans)
	}
}

/*
	TESTING THE INTEGRATION OF RAFT AND THE API
 */

func TestPostDataToLeaderAndRetrieveOnFollowerShouldPass(t *testing.T){
	//post data
	gesucht := "5"
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "store",
		"key" : "voltage",
		"value" : gesucht,
	}
	ans := SendBsonMessage(leaderaddr+":"+strconv.Itoa(cfg.MetaApiPort), msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Printf("values: %v", answerVals.Value)
	if answerVals.Value["type"] != "ok" {
		t.Errorf("test failed: %s", ans)
	}

	getmsg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "get",
		"key" : "voltage",
	}
	time.Sleep(1*time.Second)
	followerAns := SendBsonMessage(follower1+":"+strconv.Itoa(cfg.MetaApiPort), getmsg)
	followerAnswerVals := answerFormat{}
	bson.Unmarshal(followerAns, &followerAnswerVals)
	if followerAnswerVals.Value["type"] != "data" {
		t.Errorf("test failed: %s", followerAns)
	} else if followerAnswerVals.Value["value"] != gesucht{
		t.Errorf("test value not correct: %s", followerAns)

	}
}
func TestPostDataToFollowerAndRetrieveOnLeaderShouldPass(t *testing.T){
	//post data
	gesucht := "10"
	msg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "store",
		"key" : "ampere",
		"value" : gesucht,
	}
	ans := SendBsonMessage(follower1+":"+strconv.Itoa(cfg.MetaApiPort), msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	log.Printf("values: %v",answerVals.Value)
	if answerVals.Value["type"] != "ok" {
		t.Errorf("test failed: %s", ans)
	}

	getmsg := bson.M{
		"service": "electricity",
		"ip" : "1.2.3.4",
		"type" : "get",
		"key" : "ampere",
	}
	time.Sleep(1*time.Second)
	followerAns := SendBsonMessage(leaderaddr+":"+strconv.Itoa(cfg.MetaApiPort), getmsg)
	followerAnswerVals := answerFormat{}
	bson.Unmarshal(followerAns, &followerAnswerVals)
	if followerAnswerVals.Value["type"] != "data" {
		t.Errorf("test failed: %s", followerAns)
	} else if followerAnswerVals.Value["value"] != gesucht{
		t.Errorf("test value not correct: %s", followerAns)

	}
}