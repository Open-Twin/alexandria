package main

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
)

func main() {
	fmt.Println("App started")

	communication.CheckNode()
}
