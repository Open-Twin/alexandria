package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
)

func main() {
	conf := cfg.ReadConf()

	fmt.Println(conf.Welcome)
	fmt.Println(conf.ExampleTimeout)
	fmt.Println(conf.EndpointPort)
}
