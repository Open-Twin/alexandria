package main

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/Open-Twin/alexandria/raft"
	"log"
	"os"
)

func main() {
	//init raft
	//read config
	conf := cfg.ReadConf()
	log.Print("Config: ")
	log.Print(conf)
	raftLogger := log.New(os.Stdout,"raft: ",log.Ltime)

	raftNode, err := raft.Start(&conf, raftLogger)

	if err != nil{
		raftLogger.Println("error creating node. EXITING")
		os.Exit(1)
	}

	//TODO: race conditions locks???
	//dns entrypoint
	dnsEntrypointLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node: raftNode,
		Address: conf.HttpAddr,
		Logger: &dnsEntrypointLogger,
	}
	raftLogger.Println("Starting DNS entrypoint")
	dnsEntrypoint.Start()

	//dns api
	apiLogger := *log.New(os.Stdout, "dns: ", log.Ltime)
	api := &communication.API{
		Node: raftNode,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress: conf.DnsApiAddr,
		NetworkType: "udp",
		Logger:      &apiLogger,
	}
	raftLogger.Println("Starting storage API")
	api.Start()

	healthchecks := loadbalancing.HealthCheck{
		Nodes:     &raftNode.Fsm.DnsRepo.LbInfo,
		Interval:  30 * 1000,
		CheckType: loadbalancing.PingCheck,
	}
	raftLogger.Println("Starting healthchecks")
	healthchecks.ScheduleHealthChecks()

	loadbalancing.StartLoadReporting()

	httpLogger := *log.New(os.Stdout, "http: ", log.Ltime)
	service := &communication.HttpServer{
		Node:    raftNode,
		Address: conf.HttpAddr,
		UdpPort: conf.UdpPort,
		Logger:  &httpLogger,
	}
	//starts the http service (not in a goroutine so it blocks from exiting)
	raftLogger.Println("Starting HTTP service")
	service.Start()
}
