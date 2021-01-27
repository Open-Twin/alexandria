package raft

import (
	"fmt"
	"github.com/Open-Twin/alexandria/raft/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Main(){
	//read config
	rawConfig := config.ReadRawConfig()
	//validate config
	conf, err := rawConfig.ValidateConfig()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration errors - %s\n", err)
		os.Exit(1)
	}
	raftLogger := log.New(os.Stdout,"raft: ",log.Ltime)
	raftNode, err2 := NewNode(conf, raftLogger)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error configuring Node: %s", err2)
		os.Exit(1)
	}

	//attempts to join a Node if join Address is given
	if conf.JoinAddress != "" {
		go func() {
			retryJoin := func() error {
				joinUrl := url.URL{
					Scheme: "http",
					Host:   conf.JoinAddress,
					Path:   "join",
				}

				req, err := http.NewRequest(http.MethodPost, joinUrl.String(), nil)
				if err != nil {
					return err
				}
				req.Header.Add("Peer-Address", conf.RaftAddress.String())

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
					raftLogger.Print("error joining cluster")
					//Logger.Error().Err(err).Str("component", "join").Msg("Error joining cluster")
					time.Sleep(1 * time.Second)
				} else {
					break
				}
			}
		}()
	}

	httpLogger := *log.New(os.Stdout,"http: ",log.Ltime)
	service := &HttpServer{
		Node:    raftNode,
		Address: conf.HTTPAddress,
		Logger:  &httpLogger,
	}
	//starts the http service
	service.Start()
}

