package main

import (
	"github.com/Open-Twin/alexandria/loadbalancing"
)

func main() {
	loadbalancer := loadbalancing.AlexandriaBalancer{
		DnsPort: 53,
	}
	loadbalancer.StartAlexandriaLoadbalancer()
}
