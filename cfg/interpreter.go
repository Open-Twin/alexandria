package cfg

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	WELCOME         = "WELCOME"
	EXAMPLE_TIMEOUT = "EXAMPLE_TIMEOUT"
	ENDPOINT_PORT   = "ENDPOINT_PORT"

	ERROR_MSG = "The variable %s is not set."
)

type Config struct {
	Welcome        string
	ExampleTimeout int64
	EndpointPort   int64
}

func ReadConf() Config {
	fmt.Println("Reading config started")

	cfg := Config{}

	welcome := os.Getenv(WELCOME)
	if welcome == "" {
		log.Fatalf(ERROR_MSG, WELCOME)
	}
	cfg.Welcome = welcome

	exampleTimeout, err := strconv.ParseInt(os.Getenv(EXAMPLE_TIMEOUT), 10, 64)
	if err != nil {
		log.Fatalf(ERROR_MSG, EXAMPLE_TIMEOUT)
	}
	cfg.ExampleTimeout = exampleTimeout

	endpointPort, err := strconv.ParseInt(os.Getenv(ENDPOINT_PORT), 10, 64)
	if err != nil {
		log.Fatalf(ERROR_MSG, ENDPOINT_PORT)
	}
	cfg.EndpointPort = endpointPort

	fmt.Println("Reading config finished")
	return cfg
}
