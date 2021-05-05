package main

import (
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	//init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	loadbalancer := loadbalancing.AlexandriaBalancer{
		DnsPort:             53,
		HealthCheckInterval: 7 * 1000,
	}
	loadbalancer.StartAlexandriaLoadbalancer()
}
