package main

import (
	"fmt"
	"os"
	"github.com/Open-Twin/alexandria/loadbalancing"
)

func main() {
	args := os.Args[1:]

	quit := make(chan struct{})
	lb := loadbalancing.StartAlexandriaLoadbalancer(53)

	for _, ip := range args {
		lb.AddDns(ip)
	}

	fmt.Println(lb.GetDnsEntries())
	<-quit
}
