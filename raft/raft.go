package raft

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/hashicorp/raft"
	bolt "github.com/hashicorp/raft-boltdb"
	"github.com/rs/zerolog/log"
	"io"
	//stdlog "log"
	"net"
	"path/filepath"
	"time"
)



type Node struct {
	Config   *cfg.Config
	RaftNode *raft.Raft
	Fsm      *Fsm
}

/*
creates and returns a new node
*/
func NewNode(config *cfg.Config) (*Node, error){
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(config.RaftAddr.String())
	//TODO: logger
	//appLogger := log.With().Str("component", "raft-node").Logger()
	//raftConfig.Logger = stdlog.New(log.With().Str("component", "raft-node").Logger(),"",0)

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
	snapshotStoreLogger := log.With().Str("component", "raft-snapshots").Logger()
	snapshotStore, err := raft.NewFileSnapshotStore(config.DataDir,1,snapshotStoreLogger)
	if err != nil {
		return nil, err
	}
	transportLogger := log.With().Str("component", "raft-transport").Logger()
	transport, err := newTransport(config, transportLogger)
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
		log.Info().Msg("bootstrapping cluster")
	}
	return &Node{
		Config:   config,
		RaftNode: raftNode,
		Fsm:      fsm,
	}, nil
}

/*
Creates a new node but without persistent storage
only for tests
 */
func NewInMemNodeForTesting(config *cfg.Config) (*Node, error){

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
		log.Info().Msg("bootstrapping cluster")
	}
	return &Node{
		Config:   config,
		RaftNode: raftNode,
		Fsm:      fsm,
	}, nil
}
/*
creates a new tcp transport for raft
 */
func newTransport(config *cfg.Config, logger io.Writer) (*raft.NetworkTransport, error){
	address, err := net.ResolveTCPAddr("tcp",config.RaftAddr.String())
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), nil, 3, 10*time.Second, logger)

	if err != nil {
		return nil, err
	}
	return transport, nil
}