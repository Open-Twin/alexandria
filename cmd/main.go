package main

import (
	"fmt"
	"os"
	"github.com/Open-Twin/alexandria/loadbalancing"
)

func main() {
	args := os.Args[1:]

	lb := loadbalancing.StartAlexandriaLoadbalancer(1212)

	for _, ip := range args {
		lb.AddDns(ip)
	}

	fmt.Println(lb.GetDnsEntries())
}
