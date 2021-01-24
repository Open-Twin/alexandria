package raft

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"log"
	"net"
	"net/http"
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
		NewValue int `json:"newValue"`
	}{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.Logger.Println("Bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &event{
		Type:  "set",
		Value: request.NewValue,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		server.Logger.Println("")
	}
	applyFuture := server.Node.raftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Got Post: "+strconv.Itoa(request.NewValue)))
}

/*
function for handling get requests
 */
func (server *HttpServer) handleKeyGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := struct {
		Value int `json:"value"`
	}{
		Value: server.Node.fsm.stateValue,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		server.Logger.Println("")
	}

	w.Write(responseBytes)
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
