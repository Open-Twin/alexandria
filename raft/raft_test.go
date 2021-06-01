package raft_test

import (
	"bufio"
	"encoding/json"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	raftlib "github.com/hashicorp/raft"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)
const leaderaddr = "leader"
const follower1 = "follower1"
const raftaddrip = "127.0.0.1"
const raftaddrport = 7000
const httpaddrip = "127.0.0.1"
const httpaddrport = 8000
const broadcastport = 9000
const metaapiport = 20000

var s communication.HttpServer

/*
Entrypoint for the tests
*/
func TestMain(m *testing.M) {

	raftaddr := net.TCPAddr{
		IP: net.ParseIP(raftaddrip),
		Port: raftaddrport,
	}

	httpaddr := net.TCPAddr{
		IP: net.ParseIP(httpaddrip),
		Port: httpaddrport,
	}

	/*joinaddr := &net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 8000,
	}*/

	conf := cfg.Config{
		Hostname: "adin carik",
		LogLevel: 1,
		DataDir: "raft/test",
		Bootstrap: true,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: raftaddr,
		HttpAddr: httpaddr,
	}
	zerolog.SetGlobalLevel(zerolog.Level(3))
	/*followerConf := cfg.Config{
		Hostname: "follower",
		LogLevel: 1,
		DataDir: "raft/test",
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: raftaddr,
		HttpAddr: httpaddr,
		JoinAddr: joinaddr,
	}*/

	node, err := raft.NewInMemNodeForTesting(&conf)
	if err != nil{
		log.Fatal().Msg("Preparing tests failed: "+err.Error())
	}

	/*followerNode, err := raft.NewInMemNodeForTesting(&followerConf)
	if err != nil{
		log.Fatal("Preparing tests failed: "+err.Error())
	}*/

	s = communication.HttpServer{
		Node: node,
		Address: httpaddr,
		UdpPort: broadcastport,
	}
	/*followerServer := communication.HttpServer{
		Node: followerNode,
		Address: httpaddr,
	}*/
	go s.Start()
	//go followerServer.Start()
	time.Sleep(5 * time.Second)
	//checks if table exists. if not, creates it
	//runs tests
	code := m.Run()
	//cleanup
	//exits tests
	os.Exit(code)
}

/* REST API Testing */
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)
	return rr
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


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
type response struct {
	Service string `bson:"service"`
	Type string `bson:"type"`
	Key string `bson:"key"`
	Value map[string]string `bson:"value"`
}

func getMessageType(j string) string {
	var resp response
	json.Unmarshal([]byte(j), &resp)
	return resp.Value["Type"]
}
func getMessageValue(j string) string {
	var resp response
	json.Unmarshal([]byte(j), &resp)
	return resp.Value["Value"]
}

/*
JOIN requests
 */
func TestJoinWithCorrectAddressShouldPass(t *testing.T){
	//reset()

	raftAddress := raftaddrip+":"+strconv.Itoa(raftaddrport)
	req, _ := http.NewRequest("POST","/join", nil)
	req.Header.Add("Peer-Address", raftAddress)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}
/*func TestJoinWithIncorrectAddressShouldFail(t *testing.T){
	req, _ := http.NewRequest("POST","/join", nil)
	req.Header.Add("Peer-Address", "1.2.3.a")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}*/


func TestBootstrapCluster(t *testing.T){
	if s.Node.RaftNode.State() != raftlib.Leader {
		t.Errorf("bootstrapping cluster failed")
	}
}

func TestAutojoin(t *testing.T){
	//broadcast udp to find available servers

	pc, err := net.ListenPacket("udp4","")
	if err != nil {
		t.Errorf("failed autojoin listening to packet: %s", err.Error())
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", raftaddrip+":"+strconv.Itoa(broadcastport))
	if err != nil {
		t.Errorf("failed autojoin")	}
	//broadcast udp message
	_, err = pc.WriteTo([]byte("autojoin-request"), addr)
	if err != nil {
		t.Errorf("failed autojoin")	}
	//read responses

	for i:=0; i<5; i++{
		buf := make([]byte, 1024)
		pc.SetDeadline(time.Now().Add(2 * time.Second))
		n, respaddr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Info().Msgf("error reading response from \"%s: %s\n", respaddr, buf[:n])
			continue
		}
		log.Info().Msgf("autojoin received response from \"%s: %s\n", respaddr, buf[:n])
		if string(buf[:n]) == "AJ APPROVE" {
			return
		}
	}
	t.Errorf("autojoining failed because of no response")
}

/*
 * Integration Tests
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
	answerVals := response{}
	bson.Unmarshal(ans, &answerVals)
	log.Info().Msgf("values: %v",answerVals.Value)
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
	followerAnswerVals := response{}
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
	answerVals := response{}
	bson.Unmarshal(ans, &answerVals)
	log.Info().Msgf("values: %v",answerVals.Value)
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
	followerAnswerVals := response{}
	bson.Unmarshal(followerAns, &followerAnswerVals)
	if followerAnswerVals.Value["type"] != "data" {
		t.Errorf("test failed: %s", followerAns)
	} else if followerAnswerVals.Value["value"] != gesucht{
		t.Errorf("test value not correct: %s", followerAns)

	}
}