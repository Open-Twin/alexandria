package cfg

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type rawConfig struct {
	Hostname            string `validate:"required,hostname"`
	LogLevel            int    `validate:"required,max=5,min=1"`
	DataDir             string `validate:"required,dir"`
	Bootstrap           bool
	Autojoin            bool
	HealthcheckInterval int    `validate:"required,min=1000"`
	HttpAddr            string `validate:"required,ipv4"`
	RaftAddr            string `validate:"required,ipv4"`
	JoinAddr            string `validate:"omitempty,ipv4"`
	MetaApiAddr         string `validate:"required,ipv4"`
	DnsApiAddr          string `validate:"required,ipv4"`
	DnsAddr             string `validate:"required,ipv4"`
	HttpPort            int    `validate:"required,max=65536,min=1" default:"5"`
	RaftPort            int    `validate:"required,max=65536,min=1"`
	MetaApiPort         int    `validate:"required,max=65536,min=1"`
	DnsApiPort          int    `validate:"required,max=65536,min=1"`
	UdpPort             int    `validate:"required,max=65536,min=1"`
	DnsPort             int    `validate:"required,max=65536,min=1"`
}

type Config struct {
	Hostname            string
	LogLevel            int
	DataDir             string
	Bootstrap           bool
	Autojoin            bool
	HealthcheckInterval int
	HttpAddr            net.TCPAddr
	RaftAddr            net.TCPAddr
	JoinAddr            net.Addr
	MetaApiAddr         net.TCPAddr
	DnsApiAddr          net.TCPAddr
	DnsAddr             net.TCPAddr
	UdpPort             int
}

// Reads the configuration from the environment variables.
// Returns a struct with the fetched values
func ReadConf() Config {
	// This constant saves the names of the environment variables
	const (
		HOSTNAME             = "HOSTNAME"
		LOG_LEVEL            = "LOG_LEVEL"
		DATA_DIR             = "DATA_DIR"
		BOOTSTRAP            = "BOOTSTRAP"
		AUTO_JOIN            = "AUTOJOIN"
		HEALTHCHECK_INTERVAL = "HEALTHCHECK_INTERVAL"
		HTTP_ADDR            = "HTTP_ADDR"
		RAFT_ADDR            = "RAFT_ADDR"
		JOIN_ADDR            = "JOIN_ADDR"
		META_ADDR            = "META_API_ADDR"
		DNS_API_ADDR         = "DNS_API_ADDR"
		DNS_ADDR             = "DNS_ADDR"
		HTTP_PORT            = "HTTP_PORT"
		RAFT_PORT            = "RAFT_PORT"
		DNS_API_PORT         = "DNS_API_PORT"
		META_API_PORT        = "META_API_PORT"
		UDP_PORT             = "UDP_PORT"
		DNS_PORT             = "DNS_PORT"
	)

	cfg := rawConfig{}

	cfg.Hostname = os.Getenv(HOSTNAME)

	logLevel, errLog := strconv.Atoi(os.Getenv(LOG_LEVEL))
	if errLog != nil {
		logLevel = -1
	}
	cfg.LogLevel = logLevel

	cfg.DataDir = os.Getenv(DATA_DIR)

	bootStrap, errBoot := strconv.ParseBool(os.Getenv(BOOTSTRAP))
	if errBoot != nil {
		log.Warn().Msgf("Using default value for %s istead: %v\n", "Bootstrap", Bootstrap)
		bootStrap = false
	}
	cfg.Bootstrap = bootStrap

	autoJoin, errAuto := strconv.ParseBool(os.Getenv(AUTO_JOIN))
	if errAuto != nil {
		autoJoin = false
	}
	cfg.Autojoin = autoJoin

	healthInterval, errHealth := strconv.Atoi(os.Getenv(HEALTHCHECK_INTERVAL))
	if errHealth != nil {
		healthInterval = -1
	}
	cfg.HealthcheckInterval = healthInterval

	cfg.HttpAddr = os.Getenv(HTTP_ADDR)

	cfg.RaftAddr = os.Getenv(RAFT_ADDR)

	cfg.JoinAddr = os.Getenv(JOIN_ADDR)

	cfg.MetaApiAddr = os.Getenv(META_ADDR)

	cfg.DnsApiAddr = os.Getenv(DNS_API_ADDR)

	cfg.DnsAddr = os.Getenv(DNS_ADDR)

	httpPort, errHttp := strconv.Atoi(os.Getenv(HTTP_PORT))
	if errHttp != nil {
		httpPort = -1
	}
	cfg.HttpPort = httpPort

	raftPort, errPort := strconv.Atoi(os.Getenv(RAFT_PORT))
	if errPort != nil {
		raftPort = -1
	}
	cfg.RaftPort = raftPort

	metaPort, errPort := strconv.Atoi(os.Getenv(META_API_PORT))
	if errPort != nil {
		metaPort = -1
	}
	cfg.MetaApiPort = metaPort

	dnsApiPort, errPort := strconv.Atoi(os.Getenv(DNS_API_PORT))
	if errPort != nil {
		dnsApiPort = -1
	}
	cfg.DnsApiPort = dnsApiPort

	udpPort, errUdp := strconv.Atoi(os.Getenv(UDP_PORT))
	if errUdp != nil {
		udpPort = -1
	}
	cfg.UdpPort = udpPort

	dnsPort, errUdp := strconv.Atoi(os.Getenv(DNS_PORT))
	if errUdp != nil {
		dnsPort = -1
	}
	cfg.DnsPort = dnsPort

	validatedCfg, errs := validateConfig(cfg)
	for err := range errs {
		log.Error().Msgf("Error in Config: %v", err)
	}

	return validatedCfg
}

func validateConfig(rawConfig rawConfig) (Config, []validator.FieldError) {
	//Errors array in which all errors get saved
	errors := make([]validator.FieldError, 0)
	errors = nil
	//playground validator
	v := validator.New()
	err := v.Struct(rawConfig)
	//check for errors
	if err != nil {
		//loop through errors
		for _, fieldErr := range err.(validator.ValidationErrors) {
			if fieldErr.Tag() == "dir" {
				direrr := createDirectory(rawConfig.DataDir)
				if direrr == nil {
					continue
				}
			}
			setDefaultValue(fieldErr, &rawConfig)

			errors = append(errors, fieldErr)
		}
	}
	//parse ip Address
	bindAddr := net.ParseIP(rawConfig.RaftAddr)
	//create new tcpaddr from bindaddr and raftport
	raftAddr := net.TCPAddr{
		IP:   bindAddr,
		Port: rawConfig.RaftPort,
	}
	//create new tcpaddr from bindAddr and httpport
	httpAddr := net.TCPAddr{
		IP:   bindAddr,
		Port: rawConfig.HttpPort,
	}
	metaAddress := net.ParseIP(rawConfig.MetaApiAddr)
	metaAddr := net.TCPAddr{
		IP:   metaAddress,
		Port: rawConfig.MetaApiPort,
	}
	dnsApiAddress := net.ParseIP(rawConfig.DnsApiAddr)
	dnsApiAddr := net.TCPAddr{
		IP:   dnsApiAddress,
		Port: rawConfig.DnsApiPort,
	}
	dnsAddress := net.ParseIP(rawConfig.DnsAddr)
	dnsAddr := net.TCPAddr{
		IP:   dnsAddress,
		Port: rawConfig.DnsPort,
	}
	//join address
	var joinAddr net.Addr
	joinAddr = nil
	if rawConfig.JoinAddr != "" {
		joinAddress := net.ParseIP(rawConfig.JoinAddr)
		joinAddr = &net.TCPAddr{
			IP:   joinAddress,
			Port: rawConfig.HttpPort,
		}
	}
	//create config
	config := Config{
		Hostname:            rawConfig.Hostname,
		LogLevel:            rawConfig.LogLevel,
		DataDir:             rawConfig.DataDir,
		Bootstrap:           rawConfig.Bootstrap,
		Autojoin:            rawConfig.Autojoin,
		HealthcheckInterval: rawConfig.HealthcheckInterval,
		HttpAddr:            httpAddr,
		RaftAddr:            raftAddr,
		MetaApiAddr:         metaAddr,
		DnsApiAddr:          dnsApiAddr,
		DnsAddr:             dnsAddr,
		JoinAddr:            joinAddr,
		UdpPort:             rawConfig.UdpPort,
	}

	//return config
	return config, errors
}

const (
	Hostname            = "ariel"
	LogLevel            = 1
	DataDir             = "alexandria-data"
	Bootstrap           = false
	AutoJoin            = true
	HealthcheckInterval = 3000
	HttpAddr            = "127.0.0.1"
	RaftAddr            = "127.0.0.1"
	MetaApiAddr         = "0.0.0.0"
	DnsAddr             = "0.0.0.0"
	DnsApiAddr          = "0.0.0.0"
	JoinAddr            = ""
	HttpPort            = 8000
	RaftPort            = 7000
	MetaApiPort         = 20000
	DnsApiPort          = 10000
	DnsPort             = 53
	UdpPort             = 9000
)

func setDefaultValue(error validator.FieldError, conf *rawConfig) {
	log.Error().Msgf("Setting field %s threw error: %s\n", error.Field(), error.Error())

	switch error.Field() {
	case "Hostname":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "Hostname", Hostname)
		conf.Hostname = Hostname
	case "Loglevel":
		log.Warn().Msgf("Using default value for %s istead: %d\n", "Loglevel", LogLevel)
		conf.LogLevel = LogLevel
	case "DataDir":
		_, b, _, _ := runtime.Caller(0)
		path := filepath.Dir(b)
		dataDir := path + "/../" + DataDir
		log.Warn().Msgf("Using default value for %s istead: %s\n", "DataDir", dataDir)
		conf.DataDir = dataDir
		err := createDirectory(conf.DataDir)
		if err != nil {
			//TODO: fatal
			log.Fatal().Msgf("Default directory %s could not created: %s\n", conf.DataDir, err.Error())
		}
	case "Autojoin":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "Autojoin", AutoJoin)
		conf.Autojoin = AutoJoin
	case "HealthcheckInterval":
		log.Warn().Msgf("Using default value for %s istead: %d\n", "HealthcheckInterval", HealthcheckInterval)
		conf.HealthcheckInterval = HealthcheckInterval
	case "HttpAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "HttpAddr", HttpAddr)
		conf.HttpAddr = HttpAddr
	case "RaftAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "RaftAddr", RaftAddr)
		conf.RaftAddr = RaftAddr
	case "MetaApiAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "MetaApiAddr", MetaApiAddr)
		conf.MetaApiAddr = MetaApiAddr
	case "DnsApiAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "DnsApiAddr", DnsApiAddr)
		conf.DnsApiAddr = DnsApiAddr
	case "DnsAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "DnsAddr", DnsAddr)
		conf.DnsAddr = DnsAddr
	case "JoinAddr":
		log.Warn().Msgf("Using default value for %s istead: %s\n", "JoinAddr", JoinAddr)
		conf.JoinAddr = JoinAddr
	case "HttpPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "HttpPort", HttpPort)
		conf.HttpPort = HttpPort
	case "RaftPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "RaftPort", RaftPort)
		conf.RaftPort = RaftPort
	case "MetaApiPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "MetaApiPort", MetaApiPort)
		conf.MetaApiPort = MetaApiPort
	case "DnsApiPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "DnsApiPort", DnsApiPort)
		conf.DnsApiPort = DnsApiPort
	case "DnsPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "DnsPort", DnsPort)
		conf.DnsPort = DnsPort
	case "UdpPort":
		log.Warn().Msgf("Using default value for %s istead: %v\n", "UdpPort", UdpPort)
		conf.UdpPort = UdpPort
	}
}

func createDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Warn().Msgf("Directory %s not found. Creating new directory...\n", dir)
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
