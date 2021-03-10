package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/raft/config"
	"log"
	"os"
)

func main() {
	//init raft
	//read config
	rawConfig := config.ReadRawConfig()
	//validate config
	conf, confErr := rawConfig.ValidateConfig()
	if confErr != nil {
		fmt.Fprintf(os.Stderr, "Configuration errors - %s\n", confErr)
		os.Exit(1)
	}
	raftLogger := log.New(os.Stdout, "raft: ", log.Ltime)
	raftNode, err := raft.Start(conf, raftLogger)
	if err != nil {
		os.Exit(1)
	}

	//TODO: race conditions locks???
	//dns entrypoint
	dnsEntrypointLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node:    raftNode,
		Address: conf.HTTPAddress,
		Logger:  &dnsEntrypointLogger,
	}
	dnsEntrypoint.StartDnsEntrypoint()

	//dns api
	apiLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	api := &communication.API{
		Node: raftNode,
		//TODO: address and type from config
		Address:     conf.HTTPAddress,
		NetworkType: "udp",
		Logger:      &apiLogger,
	}
	api.Start()

	httpLogger := *log.New(os.Stdout, "http: ", log.Ltime)
	service := &communication.HttpServer{
		Node:    raftNode,
		Address: conf.HTTPAddress,
		Logger:  &httpLogger,
	}
	//starts the http service (not in a goroutine so it blocks from exiting)
	service.Start()

	healthchecks := loadbalancing.HealthCheck{
		Node:      raftNode,
		Interval:  5000,
		CheckType: loadbalancing.PingCheck,
	}
	healthchecks.ScheduleHealthChecks()
}
