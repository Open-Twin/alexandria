package main

import (
	"github.com/Open-Twin/alexandria/loadbalancing"
)

func main() {
	loadbalancer := loadbalancing.AlexandriaBalancer{
		DnsPort:             53,
		HealthCheckInterval: 30 * 1000,
	}
	loadbalancer.StartAlexandriaLoadbalancer()
}
