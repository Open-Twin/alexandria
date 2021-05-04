package raft_test

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

const raftaddrip = "127.0.0.1"
const raftaddrport = 7000
const httpaddrip = "127.0.0.1"
const httpaddrport = 8000

var s communication.HttpServer

/*
Entrypoint for the tests
*/
func TestMain(m *testing.M) {
	logger := log.New(os.Stdout,"",log.Ltime)

	raftaddr := net.TCPAddr{
		IP: net.ParseIP(raftaddrip),
		Port: raftaddrport,
	}

	httpaddr := net.TCPAddr{
		IP: net.ParseIP(httpaddrip),
		Port: httpaddrport,
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

	node, err := raft.NewInMemNodeForTesting(&conf)
	if err != nil{
		log.Fatal("Preparing tests failed: "+err.Error())
	}

	s = communication.HttpServer{
		Node: node,
		Address: httpaddr,
		Logger: logger,
	}
	go s.Start()
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


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
type response struct {
	Service string
	Type string
	Key string
	Value map[string]string
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
