package raft_test

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/raft/config"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"testing"
)

var s raft.HttpServer

/*
Entrypoint for the tests
*/
func TestMain(m *testing.M) {
	logger := log.New(os.Stdout,"",log.Ltime)

	raftaddr := &net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 7000,
	}
	httpaddr := &net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 8000,
	}
	conf := &config.Config{
		RaftAddress: raftaddr,
		HTTPAddress: httpaddr,
		JoinAddress: "127.0.0.1:8000",
		DataDir: "./test",
		Bootstrap: true,
	}
	node, err := raft.NewInMemNodeForTesting(conf, logger)
	if err != nil{
		log.Fatal("Preparing tests failed: "+err.Error())
	}
	s = raft.HttpServer{
		Node: node,
		Address: httpaddr,
		Logger: logger,
	}
	go s.Start()
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
type Message struct {
	Message string
	Status  bool
}
func getMessageStatus(j string) bool {
	var message Message
	json.Unmarshal([]byte(j), &message)
	return message.Status
}
func getMessage(j string) string {
	var message Message
	json.Unmarshal([]byte(j), &message)
	return message.Message
}

//NEW TEST WITH POST FOLLOWER SHOULD FAIL

func TestPostNewDataShouldPass(t *testing.T) {
	//reset()
	/*values := map[string]string{"newValue": "90"}
	json_data, _ := json.Marshal(values)
	*/
	//req, _ := http.NewRequest("POST", "/key", bytes.NewBuffer(json_data))
	req, _ := http.NewRequest("POST", "/key", strings.NewReader(`{"newValue":99}`))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); getMessageStatus(body) != true{
		t.Errorf("Expected successful POST. Got %s", body)
	}
}

func TestPostMalformedDataShouldFail(t *testing.T){
	//reset()

	req, _ := http.NewRequest("POST", "/key", strings.NewReader("newValue=falscherdatentyp"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); getMessageStatus(body) != false{
		t.Errorf("Expected failed POST due to wrong type. Got %s", body)
	}
}

func TestGetDataShouldPass(t *testing.T){
	//reset()

	req, _ := http.NewRequest("GET", "/key", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); getMessageStatus(body) != true{
		t.Errorf("Expected successful GET. Got %s", body)
	}
}

func TestJoinWithCorrectAddressShouldPass(t *testing.T){
	//reset()

	raftAddress := "localhost:8000"
	req, _ := http.NewRequest("POST","/join", nil)
	req.Header.Add("Peer-Address", raftAddress)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); getMessageStatus(body) != true{
		t.Errorf("Expected successful JOIN. Got %s", body)
	}
}
