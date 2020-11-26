package cfg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	HTTP_ADDR     = "HTTP_ADDR"
	RAFT_ADDR     = "RAFT_ADDR"
	HTTP_PORT     = "HTTP_PORT"
	RAFT_PORT     = "RAFT_PORT"
	RAFT_DATA_DIR = "RAFT_DATA_DIR"
	BOOTSTRAP     = "BOOTSTRAP"

	ERROR_MSG = "The variable %s is not set."
)

type Config struct {
	HttpAddr    string
	RaftAddr    string
	HttpPort    int64
	RaftPort    int64
	RaftDataDir string
	Bootstrap   bool
}

func ReadConf() Config {
	fmt.Println("Reading config started")

	cfg := Config{}

	http_addr := os.Getenv(HTTP_ADDR)
	if http_addr == "" {
		log.Fatalf(ERROR_MSG, HTTP_ADDR)
	}
	cfg.HttpAddr = http_addr

	raft_addr := os.Getenv(RAFT_ADDR)
	if raft_addr == "" {
		log.Fatalf(ERROR_MSG, RAFT_ADDR)
	}
	cfg.RaftAddr = raft_addr

	http_port, err := strconv.ParseInt(os.Getenv(HTTP_PORT), 10, 64)
	if err != nil {
		log.Fatalf(ERROR_MSG, HTTP_PORT)
	}
	cfg.HttpPort = http_port

	raft_port, err := strconv.ParseInt(os.Getenv(RAFT_PORT), 10, 64)
	if err != nil {
		log.Fatalf(ERROR_MSG, RAFT_PORT)
	}
	cfg.RaftPort = raft_port

	raft_data_dir := os.Getenv(RAFT_DATA_DIR)
	if raft_data_dir == "" {
		log.Fatalf(ERROR_MSG, RAFT_DATA_DIR)
	}
	cfg.RaftDataDir = raft_data_dir

	bootstrap := os.Getenv(BOOTSTRAP)
	bootstrap = strings.ToLower(bootstrap)
	if bootstrap == "true" {
		cfg.Bootstrap = true
	} else if bootstrap != "false" {
		log.Fatalf(ERROR_MSG, BOOTSTRAP)
	}
	cfg.Bootstrap = false

	fmt.Println("Reading config finished")
	return cfg
}
