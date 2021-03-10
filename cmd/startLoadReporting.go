package main

import (
	"github.com/Open-Twin/alexandria/loadbalancing"
)

func main() {
	quit := make(chan struct{})
	loadbalancing.StartLoadReporting()
	<- quit
}
