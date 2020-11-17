package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	bolt "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main(){
	rawConfig := ReadRawConfig()
	config, err := rawConfig.ValidateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration errors - %s\n", err)
		os.Exit(1)
	}

	raftNode, error := NewNode(config)
}

type node struct {
	config   *Config
	raftNode *raft.Raft
	fsm      *fsm
	log      *log.Logger
}
func NewNode(config *Config) (*node, error){

	raftConfig := raft.DefaultConfig()

	fsm := &fsm{
		stateValue : 0,
	}

	logStore, err := bolt.NewBoltStore(filepath.Join(config.DataDir,"logStore"))
	if err != nil {
		return nil, err
	}
	stableStore, err := bolt.NewBoltStore(filepath.Join(config.DataDir,"stableStore"))
	if err != nil {
		return nil, err
	}
	snapshotStore, err := raft.NewFileSnapshotStore(config.DataDir,1,)
	if err != nil {
		return nil, err
	}
	transport, err := newTransport(config)
	if err!= nil {
		return nil, err
	}
	raftNode, _ := raft.NewRaft(raftConfig,fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}
	if config.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		raftNode.BootstrapCluster(configuration)
	}
	return &node{
		config:   config,
		raftNode: raftNode,
		log:      log,
		fsm:      fsm,
	}, nil
}
func newTransport(config *Config) (*raft.NetworkTransport, error){
	address, err := net.ResolveTCPAddr("tcp",config.RaftAddress.String())
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), config.HTTPAddress, 3, 10*time.Second, log)
	if err != nil {
		return nil, err
	}
	return transport, nil
}
