package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
)

func main() {
	cfg.ReadConf()

	fmt.Println(cfg.WELCOME)
	fmt.Println(cfg.ENDPOINT_PORT)
	fmt.Println(cfg.EXAMPLE_TIMEOUT)
}
