package cfg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// This constant saves the names of the environment variables
const (
	HOSTNAME      = "HOSTNAME"
	UDP_PORT      = "UDP_PORT"
	TCP_PORT      = "TCP_PORT"
	LOG_LEVEL     = "LOG_LEVEL"
	HTTP_ADDR     = "HTTP_ADDR"
	RAFT_ADDR     = "RAFT_ADDR"
	HTTP_PORT     = "HTTP_PORT"
	RAFT_PORT     = "RAFT_PORT"
	JOIN_PORT     = "JOIN_PORT"
	RAFT_DATA_DIR = "RAFT_DATA_DIR"
	BOOTSTRAP     = "BOOTSTRAP"

	ERROR_MSG = "The variable %s is not set."
)

// Struct that saves all the configured values
type Config struct {
	Hostname    string
	UdpPort     int
	TcpPort     int
	LogLevel    int
	HttpAddr    string
	RaftAddr    string
	HttpPort    int
	RaftPort    int
	JoinPort	int
	RaftDataDir string
	Bootstrap   bool
}

// Reads the configuration from the environment variables.
// Returns a struct with the fetched values
func ReadConf() Config {
	fmt.Println("Reading config started")

	cfg := Config{}

	hostname := os.Getenv(HOSTNAME)
	if hostname == "" {
		log.Fatalf(ERROR_MSG, HOSTNAME)
	}
	cfg.Hostname = hostname

	udp_port, err := strconv.Atoi(os.Getenv(UDP_PORT))
	if err != nil {
		log.Fatalf(ERROR_MSG, HTTP_PORT)
	}
	cfg.UdpPort = udp_port

	tcp_port, err := strconv.Atoi(os.Getenv(TCP_PORT))
	if err != nil {
		log.Fatalf(ERROR_MSG, HTTP_PORT)
	}
	cfg.TcpPort = tcp_port

	log_level, err := strconv.Atoi(os.Getenv(LOG_LEVEL))
	if err != nil {
		log.Fatalf(ERROR_MSG, HTTP_PORT)
	}
	cfg.LogLevel = log_level

	http_addr := os.Getenv(HTTP_ADDR)
	/*if http_addr == "" {
		log.Fatalf(ERROR_MSG, HTTP_ADDR)
	}*/
	cfg.HttpAddr = http_addr

	raft_addr := os.Getenv(RAFT_ADDR)
	if raft_addr == "" {
		log.Fatalf(ERROR_MSG, RAFT_ADDR)
	}
	cfg.RaftAddr = raft_addr

	http_port, err := strconv.Atoi(os.Getenv(HTTP_PORT))
	if err != nil {
		log.Fatalf(ERROR_MSG, HTTP_PORT)
	}
	cfg.HttpPort = http_port

	raft_port, err := strconv.Atoi(os.Getenv(RAFT_PORT))
	if err != nil {
		log.Fatalf(ERROR_MSG, RAFT_PORT)
	}
	cfg.RaftPort = raft_port

	join_port, err := strconv.Atoi(os.Getenv(JOIN_PORT))
	if err != nil {
		log.Fatalf(ERROR_MSG, JOIN_PORT)
	}
	cfg.JoinPort = join_port

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
