package raft_test

import (
	"bytes"
	"encoding/json"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/raft/config"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var s raft.HttpServer
const raftaddrip = "127.0.0.1"
const raftaddrport = 7000
const httpaddrip = "127.0.0.1"
const httpaddrport = 8000
/*
Entrypoint for the tests
*/
func TestMain(m *testing.M) {
	logger := log.New(os.Stdout,"",log.Ltime)

	raftaddr := &net.TCPAddr{
		IP: net.ParseIP(raftaddrip),
		Port: raftaddrport,
	}
	httpaddr := &net.TCPAddr{
		IP: net.ParseIP(httpaddrip),
		Port: httpaddrport,
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

//NEW TEST WITH POST FOLLOWER SHOULD FAIL

/*
POST requests
 */
func TestPostNewDataShouldPass(t *testing.T) {
	//reset()
	values := map[string]string{
		"service": "electricity",
		"ip": "1.2.3.4",
		"type": "store",
		"key": "volt",
		"value": "50",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedType := "ok"

	if body := response.Body.String(); getMessageType(body) != expectedType{
		t.Errorf("Expected successful POST. Got %s", body)
	}
}

func TestPostDeleteDataShouldPass(t *testing.T) {
	//reset()
	values := map[string]string{
		"service": "deleteServiceTest",
		"ip": "1.2.3.4",
		"type": "store",
		"key": "volt",
		"value": "50",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	values = map[string]string{
		"service": "deleteServiceTest",
		"ip": "1.2.3.4",
		"type": "delete",
		"key": "volt",
	}
	jsonData, _ = json.Marshal(values)

	req, _ = http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedType := "ok"

	if body := response.Body.String(); getMessageType(body) != expectedType{
		t.Errorf("Expected successful delete. Got %s", body)
	}
}

func TestPostUpdateDataShouldPass(t *testing.T) {
	//reset()
	values := map[string]string{
		"service": "updateServiceTest",
		"ip": "1.2.3.4",
		"type": "store",
		"key": "volt",
		"value": "20",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	expectedValue := "500"
	values = map[string]string{
		"service": "deleteServiceTest",
		"ip": "1.2.3.4",
		"type": "update",
		"key": "volt",
		"value": expectedValue,
	}
	jsonData, _ = json.Marshal(values)

	req, _ = http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)


	values = map[string]string{
		"service": "deleteServiceTest",
		"ip": "1.2.3.4",
		"type": "get",
		"key": "volt",
	}
	jsonData, _ = json.Marshal(values)

	req, _ = http.NewRequest("GET", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); getMessageValue(body) != expectedValue{
		t.Errorf("Expected successful update. Got %s", body)
	}
}
/*
GET requests
 */
func TestGetDataShouldPass(t *testing.T){
	//reset()
	//Post data to get later
	expectedValue := "100"
	postValues := map[string]string{
		"service": "water",
		"ip": "1.2.3.4",
		"type": "store",
		"key": "height",
		"value": expectedValue,
	}
	jsonData, _ := json.Marshal(postValues)
	req, _ := http.NewRequest("POST", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	values := map[string]string{
		"service": "water",
		"ip": "1.2.3.4",
		"type": "get",
		"key": "height",
	}
	jsonData, _ = json.Marshal(values)

	req, _ = http.NewRequest("GET", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedType := "data"

	if body := response.Body.String(); getMessageType(body) != expectedType && getMessageValue(body) != "100"{
		t.Errorf("Expected successful GET. Got %s", body)
	}
}

func TestGetNonExistingShouldFail(t *testing.T){
	//reset()
	//Post data to get later

	values := map[string]string{
		"service": "nonexisting",
		"ip": "1.2.3.4",
		"type": "get",
		"key": "height",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("GET", "/key", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expectedType := "error"

	if body := response.Body.String(); getMessageType(body) != expectedType{
		t.Errorf("Expected successful GET. Got %s", body)
	}
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
