package communication

import (
	"bytes"
	"fmt"
	raft2 "github.com/Open-Twin/alexandria/raft"
	"github.com/hashicorp/raft"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HttpServer struct {
	Node    *raft2.Node
	Address net.TCPAddr
	Logger  *log.Logger
}

/*
Starts the webservice
 */
func (server *HttpServer) Start() {
	server.Logger.Printf("Start listening for auto-join requests")
	go startAutojoinListener(strings.Split(server.Address.String(),":")[1])
	server.Logger.Printf("Starting server with Address %v\n", server.Address.String())
	if err := http.ListenAndServe(server.Address.String(), server); err != nil {
		server.Logger.Fatal("Error running HTTP server")
	}
}
func startAutojoinListener(port string){
	//TODO: use udp endpoint
	log.Print("PORT:"+port)
	//TODO: autojoin port
	pc, err := net.ListenPacket("udp4", ":"+port)
	if err != nil {
		panic(err)
	}
	//pc.Close()
	for{
		buf := make([]byte, 1024)
		n,addr,err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		log.Printf("AUTOJOIN %s sent this: %s\n", addr, buf[:n])
		pc.WriteTo([]byte("AJ APPROVE"), addr)
	}
}

/*
Differentiates between /key and /join requests and forwards them to
the appropriate function
 */
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, "/join") {
		fmt.Println("JOIN REQUEST")
		server.handleJoin(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
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
		return
	}
	leaderAddr := server.Node.Config.JoinAddr
	log.Print("State: "+server.Node.RaftNode.State().String()+ " Leader addr: "+leaderAddr.String())
	if server.Node.RaftNode.State() != raft.Leader {
		server.Logger.Print("forwarding join to leader")
		forwardJoinToLeader(peerAddress, leaderAddr.String(), w)
		//sendResponse(request.Service,request.Key,"ok",,w)
		return
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
/*
func forwardToLeader(eventBytes []byte, path, leaderAddr string, request *event, w http.ResponseWriter) {
	leaderUrl := url.URL{
		Scheme: "http",
		//TODO: Leader address
		Host:   leaderAddr,
		Path:   path,
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
	log.Print("Request forwarded to leader "+leaderAddr)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(bodyBytes)
}*/

func forwardJoinToLeader(peerAddr, leaderAddr string, w http.ResponseWriter){
	leaderUrl := url.URL{
		Scheme: "http",
		//TODO: Leader address
		Host:   leaderAddr,
		Path:   "join",
	}

	req, _ := http.NewRequest("POST", leaderUrl.String(), bytes.NewBuffer([]byte(peerAddr)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		sendHttpResponse("join",peerAddr,"error",err.Error(),w)
		return
	}else if resp.StatusCode != http.StatusOK {
		sendHttpResponse("join",peerAddr,"error","non 200 status code: "+strconv.Itoa(resp.StatusCode),w)
		return
	}
	log.Print("join request forwarded to leader "+leaderAddr)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//w.Write(bodyBytes)
	sendHttpResponse("join",peerAddr,"ok",string(bodyBytes),w)
}
