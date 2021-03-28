package storageApi

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/storage"
	raftlib "github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func handleMetadata(buf []byte, node *raft.Node) []byte {
	request := struct {
		Service string `bson:"Service"`
		Ip      string `bson:"Ip"`
		Type    string `bson:"Type"`
		Key     string `bson:"Key"`
		Value   string `bson:"Value"`
	}{}

	if err := bson.Unmarshal(buf, &request); err != nil {
		log.Error().Msgf("Bad request: %v", err.Error())
		return communication.CreateResponse("", "error", "something went wrong. please check your input.")
	}
	//handle get request
	if request.Type == "get" {
		data, err := node.Fsm.MetadataRepo.Read(request.Service, request.Ip, request.Key)
		if err != nil {
			//return error
			return communication.CreateMetadataResponse(request.Service, request.Key, "error", err.Error())
		}
		return communication.CreateMetadataResponse(request.Service, request.Key, "data", data)
	}
	//handle other requests
	if node.RaftNode.State() != raftlib.Leader {
		resp, err := forwardToLeader(buf, string(node.RaftNode.Leader()), node.Config.MetaApiAddr.Port)
		if err != nil{
			log.Error().Msg("forward to leader failed")
			return communication.CreateResponse("", "error", "something went wrong. please check your input.")
		}
		return resp
	}
	//marshal record
	event := storage.Metadata{
		Dnsormetadata: false,
		Service:       request.Service,
		Ip:            request.Ip,
		Type:          request.Type,
		Key:           request.Key,
		Value:         request.Value,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error().Msg("unexpected failure")
	}

	//Apply to Raft cluster
	applyFuture := node.RaftNode.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		log.Error().Msgf("could not apply to raft cluster: %v", err.Error())
		return communication.CreateMetadataResponse(request.Service, request.Key, "error", err.Error())
	}
	var resp []byte
	if err != nil {
		resp = communication.CreateMetadataResponse(request.Service, request.Key, "error", "something went wrong. please check your input.")
	} else {
		resp = communication.CreateMetadataResponse(request.Service, request.Key, "ok", "null")
	}
	return resp
}
