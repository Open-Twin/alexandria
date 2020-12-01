package raft

import (
	"github.com/hashicorp/raft"
	bolt "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"path/filepath"
	"time"
)



type node struct {
	config   *Config
	raftNode *raft.Raft
	fsm      *fsm
	log      *log.Logger
}
func NewNode(config *Config, log *log.Logger) (*node, error){

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
	snapshotStoreLogger := log.Writer()
	snapshotStore, err := raft.NewFileSnapshotStore(config.DataDir,1,snapshotStoreLogger)
	if err != nil {
		return nil, err
	}
	transport, err := newTransport(config, log)
	if err!= nil {
		return nil, err
	}
	raftNode, err := raft.NewRaft(raftConfig,fsm, logStore, stableStore, snapshotStore, transport)
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
func newTransport(config *Config, logger *log.Logger) (*raft.NetworkTransport, error){
	address, err := net.ResolveTCPAddr("tcp",config.RaftAddress.String())
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), config.HTTPAddress, 3, 10*time.Second, logger)
	if err != nil {
		return nil, err
	}
	return transport, nil
}