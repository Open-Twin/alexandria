package raft

import (
	"encoding/json"
	"github.com/hashicorp/raft"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type httpServer struct {
	node *node
	address net.Addr
	logger *log.Logger
}

func (server *httpServer) Start() {
	server.logger.Printf("Starting server with address %v\n", server.address.String())

	if err := http.ListenAndServe(server.address.String(), nil); err != nil {
		server.logger.Fatal("Error running HTTP server")
	}
}

func (server *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/key") {
		server.handleRequest(w, r)
	} else if strings.Contains(r.URL.Path, "/join") {
		server.handleJoin(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (server *httpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
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

func (server *httpServer) handleKeyPost(w http.ResponseWriter, r *http.Request) {
	request := struct {
		NewValue int `json:"newValue"`
	}{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.logger.Println("Bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &event{
		Type:  "set",
		Value: request.NewValue,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		server.logger.Println("")
	}
	applyFuture := server.node.raftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *httpServer) handleKeyGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := struct {
		Value int `json:"value"`
	}{
		Value: server.node.fsm.stateValue,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		server.logger.Println("")
	}

	w.Write(responseBytes)
}

func (server *httpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	peerAddress := r.Header.Get("Peer-Address")
	if peerAddress == "" {
		server.logger.Println("Peer-Address not set on request")
		w.WriteHeader(http.StatusBadRequest)
	}

	addPeerFuture := server.node.raftNode.AddVoter(
		raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		server.logger.Printf("\"Error joining peer to Raft\" %v\n", peerAddress)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.logger.Printf("Peer joined Raft with address %v\n",peerAddress)
	w.WriteHeader(http.StatusOK)
}
