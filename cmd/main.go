package main

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/communication/storageApi"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	//init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	//read config
	conf := cfg.ReadConf()
	logLevel := conf.LogLevel
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	log.Debug().Msgf("Config: %v", conf)

	raftNode, err := raft.Start(&conf)
	if err != nil {
		log.Panic().Msg("Error creating node. Exiting!")
		os.Exit(1)
	}

	//TODO: race conditions locks???
	//dns entrypoint
	dnsEntrypoint := &communication.DnsEntrypoint{
		Node:    raftNode,
		Address: conf.DnsAddr,
	}
	log.Info().Msg("Starting DNS entrypoint")
	dnsEntrypoint.Start()

	//dns api
	api := &storageApi.API{
		Node: raftNode,
		//TODO: address and type from config
		MetaAddress: conf.MetaApiAddr,
		DNSAddress:  conf.DnsApiAddr,
		NetworkType: "udp",
	}
	log.Info().Msg("Starting storage API")
	api.Start()

	healthchecks := loadbalancing.HealthCheck{
		Nodes:     raftNode.Fsm.DnsRepo.LbInfo,
		Interval:  30 * 1000,
		CheckType: loadbalancing.PingCheck,
	}
	log.Info().Msg("Starting healthchecks")
	healthchecks.ScheduleHealthChecks()

	//log.Info().Msg("Starting reporting load to loadbalancer")
	//loadbalancing.StartLoadReporting()

	service := &communication.HttpServer{
		Node:    raftNode,
		Address: conf.HttpAddr,
		UdpPort: conf.UdpPort,
	}
	//starts the http service (not in a goroutine so it blocks from exiting)
	log.Info().Msg("Starting HTTP service")
	service.Start()
}
