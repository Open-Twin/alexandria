package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
)

func main() {
	cfg.ReadConf()

	fmt.Println(cfg.StringConf)
	fmt.Println(cfg.NumberConf)
	fmt.Println(cfg.StringConf["WELCOME"])
}
