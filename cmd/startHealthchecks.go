package main

import (
	"github.com/Open-Twin/alexandria/loadbalancing"
	"os"
	"strconv"
)

func main() {
	nodes := os.Args[2:]
	interval, _ := strconv.Atoi(os.Args[1])

	hc := loadbalancing.HealthCheck{}
	for _, ip := range nodes {
		hc.AddNode(ip)
	}

	quit := make(chan struct{})
	hc.ScheduleHealthChecks(interval)
	<- quit
}
