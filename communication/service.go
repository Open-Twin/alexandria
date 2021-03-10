package communication

import (
	"encoding/json"
	"fmt"
	raft2 "github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/hashicorp/raft"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type HttpServer struct {
	Node    *raft2.Node
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

	event := storage.Metadata{
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
	//Apply to Raft cluster
	applyFuture := server.Node.RaftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		server.Logger.Println("could not apply to raft cluster: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err != nil{
		sendHttpResponse(request.Service,request.Key,"error","something went wrong. please check your input.",w)
	}else {
		sendHttpResponse(request.Service,request.Key,"ok","null",w)
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

	respValue, err := server.Node.Fsm.MetadataRepo.Read(request.Service,request.Ip,request.Key)

	if err != nil{
		sendHttpResponse(request.Service,request.Key,"error",err.Error(),w)

	}else {
		sendHttpResponse(request.Service,request.Key,"data",respValue,w)
	}
}

/*
handles a /join request and attempts to join the Node
*/
func (server *HttpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	peerAddress := r.Header.Get("Peer-Address")
	if peerAddress == "" {
		server.Logger.Println("Peer-Address not set on request")
		w.WriteHeader(http.StatusBadRequest)
	}

	addPeerFuture := server.Node.RaftNode.AddVoter(
		raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		server.Logger.Printf("\"Error joining peer to Raft\" %v\n", peerAddress)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.Logger.Printf("Peer joined Raft with Address %v\n",peerAddress)
	w.WriteHeader(http.StatusOK)
}

func sendHttpResponse(service, key, etype, value string, w http.ResponseWriter){

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


