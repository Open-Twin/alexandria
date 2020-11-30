package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	bolt "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"net/http"
	"net/url"
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
	raftLogger := log.New(os.Stdout,"raft",log.Ltime)
	raftNode, err := NewNode(config, raftLogger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring node: %s", err)
		os.Exit(1)
	}

	if config.JoinAddress != "" {
		go func() {
			retryJoin := func() error {
				joinUrl := url.URL{
					Scheme: "http",
					Host:   config.JoinAddress,
					Path:   "join",
				}

				req, err := http.NewRequest(http.MethodPost, joinUrl.String(), nil)
				if err != nil {
					return err
				}
				req.Header.Add("Peer-Address", config.RaftAddress.String())

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("non 200 status code: %d", resp.StatusCode)
				}

				return nil
			}

			for {
				if err := retryJoin(); err != nil {
					//logger.Error().Err(err).Str("component", "join").Msg("Error joining cluster")
					fmt.Println("Error joining cluster")
					time.Sleep(1 * time.Second)
				} else {
					break
				}
			}
		}()
	}
	httpLogger := log.New(os.Stdout,"http",log.Ltime)
	service := &service{
		node: raftNode,
		address: config.HTTPAddress,
		logger: &httpLogger,
	}
	service.Start()
}

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