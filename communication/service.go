package communication

import (
	"bytes"
	"github.com/Open-Twin/alexandria/raft"
	raftlib "github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HttpServer struct {
	Node    *raft.Node
	Address net.TCPAddr
	UdpPort int
}

/*
Starts the webservice
 */
func (server *HttpServer) Start() {
	log.Info().Msg("Start listening for auto-join requests")
	go startAutojoinListener(server.UdpPort)
	log.Info().Msgf("Starting server with Address %v\n", server.Address.String())
	if err := http.ListenAndServe(server.Address.String(), server); err != nil {
		log.Fatal().Msgf("Error running HTTP server, Error: %v", err)
	}
}

func startAutojoinListener(port int){
	//log.Print("autojoin port:"+strconv.Itoa(port))
	pc, err := net.ListenPacket("udp4", ":"+strconv.Itoa(port))
	if err != nil {
		//TODO panic
		panic(err)
	}
	//pc.Close()
	for{
		buf := make([]byte, 1024)
		n,addr,err := pc.ReadFrom(buf)
		if err != nil {
			log.Error().Msgf("error on reading autojoin response: %s",err.Error())
		}

		log.Info().Msgf("AUTOJOIN %s sent this: %s\n", addr, buf[:n])
		pc.WriteTo([]byte("AJ APPROVE"), addr)
	}
}

/*
Differentiates between /key and /join requests and forwards them to
the appropriate function
 */
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/join") {
		log.Info().Msgf("received join request from %s", r.RemoteAddr)
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
		log.Warn().Msg("Peer-Address not set on request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//leaderAddr := server.Node.Config.JoinAddr
	leaderAddr := server.Node.RaftNode.Leader()
	log.Info().Msg("State: "+server.Node.RaftNode.State().String()+ " Leader addr: "+string(leaderAddr))
	if server.Node.RaftNode.State() != raftlib.Leader {
		log.Info().Msg("forwarding join to leader")
		forwardJoinToLeader(peerAddress, string(leaderAddr), server.Address.Port, w)
		//sendResponse(request.Service,request.Key,"ok",,w)
		return
	}

	addPeerFuture := server.Node.RaftNode.AddVoter(
		raftlib.ServerID(peerAddress), raftlib.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		log.Warn().Msgf("\"Error joining peer to Raft\" %v\n", peerAddress)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	log.Info().Msgf("Peer joined Raft with Address %v\n",peerAddress)
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

func forwardJoinToLeader(peerAddr, leaderAddr string, httpport int, w http.ResponseWriter){
	leaderUrl := url.URL{
		Scheme: "http",
		//TODO: Leader address
		Host:   strings.Split(leaderAddr,":")[0]+strconv.Itoa(httpport),
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
		log.Fatal().Msg(err.Error())
	}
	//w.Write(bodyBytes)
	sendHttpResponse("join",peerAddr,"ok",string(bodyBytes),w)
}
