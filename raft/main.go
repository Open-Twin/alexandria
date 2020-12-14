package raft

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Main(){
	rawConfig := ReadRawConfig()

	/*bind := os.Args[1]
	join := os.Args[2]
	raftport, erri := strconv.Atoi(os.Args[3])
	httpport, erri := strconv.Atoi(os.Args[4])
	joinport, erri := strconv.Atoi(os.Args[5])
	bootstrap, erri := strconv.ParseBool(os.Args[6])
	if erri != nil {

	}
	rawConfig := RawConfig{
		BindAddress: bind,
		JoinAddress: join,
		RaftPort: raftport,
		HTTPPort: httpport,
		JoinPort: joinport,
		DataDir: "./raft/test",
		Bootstrap: bootstrap,
	}*/
	config, err := rawConfig.ValidateConfig()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration errors - %s\n", err)
		os.Exit(1)
	}
	raftLogger := log.New(os.Stdout,"raft",log.Ltime)
	raftNode, err2 := NewNode(config, raftLogger)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error configuring node: %s", err2)
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
	httpLogger := *log.New(os.Stdout,"http",log.Ltime)
	service := &httpServer{
		node: raftNode,
		address: config.HTTPAddress,
		logger: &httpLogger,
	}
	service.Start()
}

