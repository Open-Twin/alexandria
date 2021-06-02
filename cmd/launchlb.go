package main

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	// init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// read configuration
	conf := cfg.ReadConf()
	logLevel := conf.LogLevel
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	log.Debug().Msgf("Config: %v", conf)

	loadbalancer := loadbalancing.AlexandriaBalancer{
		RegistrationPort:    cfg.LbPort,
		DnsPort:             cfg.DnsPort,
		DnsApiPort:          cfg.DnsApiPort,
		MetdataApiPort:      cfg.MetaApiPort,
		HttpPingPort:        cfg.HttpPingPort,
		HealthCheckInterval: cfg.HealthcheckInterval * time.Millisecond,
		RemoveNodeTimeout:   cfg.RemoveNodeTimeout * time.Second,
		RequestTimeout:      cfg.HealthcheckRequestTimeout * time.Millisecond,
	}
	loadbalancer.StartAlexandriaLoadbalancer()
}
