package raft

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/hashicorp/raft"
	bolt "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)



type Node struct {
	Config   *cfg.Config
	RaftNode *raft.Raft
	Fsm      *Fsm
	logger   *log.Logger
}

/*
creates and returns a new node
*/
func NewNode(config *cfg.Config, logger *log.Logger) (*Node, error){
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(config.RaftAddr.String())
	//raftConfig.Logger = log.New(logger, "", 0)

	metarepo := storage.NewInMemoryStorageRepository()
	dnsrepo := storage.NewInMemoryDNSStorageRepository()
	fsm := &Fsm{
		MetadataRepo: metarepo,
		DnsRepo: dnsrepo,
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
	transport, err := newTransport(config, logger)
	if err != nil {
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
		logger.Print("bootstrapping cluster")
	}
	return &Node{
		Config:   config,
		RaftNode: raftNode,
		logger:   logger,
		Fsm:      fsm,
	}, nil
}

/*
Creates a new node but without persistent storage
only for tests
 */
func NewInMemNodeForTesting(config *cfg.Config, logger *log.Logger) (*Node, error){

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(config.RaftAddr.String())
	//raftConfig.Logger = log.New(Logger, "", 0)
	metarepo := storage.NewInMemoryStorageRepository()
	dnsrepo := storage.NewInMemoryDNSStorageRepository()
	fsm := &Fsm{
		MetadataRepo: metarepo,
		DnsRepo: dnsrepo,
	}

	logStore := raft.NewInmemStore()

	stableStore := raft.NewInmemStore()

	snapshotStore := raft.NewInmemSnapshotStore()

	_, transport := raft.NewInmemTransport(raft.ServerAddress(config.RaftAddr.String()))

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
		logger.Print("bootstrapping cluster")
	}
	return &Node{
		Config:   config,
		RaftNode: raftNode,
		logger:   logger,
		Fsm:      fsm,
	}, nil
}
/*
creates a new tcp transport for raft
 */
func newTransport(config *cfg.Config, logger *log.Logger) (*raft.NetworkTransport, error){
	address, err := net.ResolveTCPAddr("tcp",config.RaftAddr.String())
	if err != nil {
		return nil, err
	}

	//TODO Logger statt stdout
	transport, err := raft.NewTCPTransport(address.String(), nil, 3, 10*time.Second, os.Stdout)

	if err != nil {
		return nil, err
	}
	return transport, nil
}