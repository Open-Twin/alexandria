package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HttpServer struct {
	Node    *node
	Address net.Addr
	Logger  *log.Logger
}

/*
Starts the webservice
 */
func (server *HttpServer) Start() {
	server.Logger.Printf("Starting server with Address %v\n", server.Address.String())

	if err := http.ListenAndServe(server.Address.String(), server); err != nil {
		server.Logger.Fatal("Error running HTTP server")
	}
}

/*
Differentiates between /key and /join requests and forwards them to
the appropriate function
 */
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, "/key") {
		fmt.Println("KEY REQUEST")
		server.handleRequest(w, r)
	} else if strings.Contains(r.URL.Path, "/join") {
		fmt.Println("JOIN REQUEST")
		server.handleJoin(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

/*
   Handles /key requests and differentiates between post and get
*/
func (server *HttpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		server.handleKeyPost(w, r)
		return
	case http.MethodGet:
		server.handleKeyGet(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

/*
function for handling post requests
 */
func (server *HttpServer) handleKeyPost(w http.ResponseWriter, r *http.Request) {
	request := struct {
		Service string
		Ip string
		Type string
		Key string
		Value string
	}{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.Logger.Println("Bad request")
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		server.Logger.Print(string(bodyBytes))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &event{
		Service: request.Service,
		Ip: request.Ip,
		Type:  request.Type,
		Key: request.Key,
		Value: request.Value,
	}
	log.Print("DEJAN2: "+event.Service+" "+event.Ip+" "+event.Type+" "+event.Key+" "+event.Value)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		server.Logger.Println("")
	}

	//TODO: forward to leader if not leader
	log.Print("State: "+server.Node.raftNode.State().String()+ " Leader addr: "+server.Node.config.JoinAddress)
	if server.Node.raftNode.State() != raft.Leader {
		leaderUrl := url.URL{
			Scheme: "http",
			//TODO: Leader address
			Host:   server.Node.config.JoinAddress,
			Path:   "key",
		}

		req, _ := http.NewRequest("POST", leaderUrl.String(), bytes.NewBuffer(eventBytes))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		req.Close = true
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			sendResponse(request.Service,request.Key,"error",err.Error(),w)
			return
		}else if resp.StatusCode != http.StatusOK {
			sendResponse(request.Service,request.Key,"error","non 200 status code: "+strconv.Itoa(resp.StatusCode),w)
			return
		}
		log.Print("Request forwarded to leader")
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(bodyBytes)
		//sendResponse(request.Service,request.Key,"ok",,w)
		return
	}

	//Apply to Raft cluster
	applyFuture := server.Node.raftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		server.Logger.Println("could not apply to raft cluster: "+err.Error())
		sendResponse(request.Service,request.Key,"error","could not apply to raft cluster: "+err.Error(),w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err != nil{
		sendResponse(request.Service,request.Key,"error","something went wrong. please check your input.",w)
	}else {
		sendResponse(request.Service,request.Key,"ok","null",w)
	}
}

/*
function for handling get requests
 */
func (server *HttpServer) handleKeyGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	request := struct {
		Service string
		Ip string
		Type string
		Key string
	}{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.Logger.Println("Bad request")
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		server.Logger.Print(string(bodyBytes))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Print("DEJAN3: "+request.Service+" "+request.Ip+" "+request.Type+" "+request.Key)

	respValue, err := server.Node.fsm.Repo.Read(request.Service,request.Ip,request.Key)

	if err != nil{
		sendResponse(request.Service,request.Key,"error",err.Error(),w)

	}else {
		sendResponse(request.Service,request.Key,"data",respValue,w)
	}
}

/*
handles a /join request and attempts to join the node
 */
func (server *HttpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	peerAddress := r.Header.Get("Peer-Address")
	if peerAddress == "" {
		server.Logger.Println("Peer-Address not set on request")
		w.WriteHeader(http.StatusBadRequest)
	}

	addPeerFuture := server.Node.raftNode.AddVoter(
		raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		server.Logger.Printf("\"Error joining peer to Raft\" %v\n", peerAddress)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.Logger.Printf("Peer joined Raft with Address %v\n",peerAddress)
	w.WriteHeader(http.StatusOK)
}

func sendResponse(service, key, etype, value string, w http.ResponseWriter){

	valueMap := map[string]string{
		"Type": etype,
		"Value": value,
	}
	response := struct {
		Service string
		Type string
		Key string
		Value map[string]string
	}{
		Service: service,
		Type: "response",
		Key: key,
		Value: valueMap,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		//server.Logger.Println("")
		log.Print("sendresponse failed")
	}

	w.Write(responseBytes)
}
