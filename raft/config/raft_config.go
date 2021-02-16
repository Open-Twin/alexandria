package config

import (
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/go-playground/validator/v10"
	"log"
	"net"
	"os"
	"strconv"
)

type RawConfig struct {
	BindAddress string `validate:"required,ipv4"`
	JoinAddress string `validate:"omitempty"`    //ipv4 not working with urls -> TODO
	HTTPAddress string
	RaftPort    int `validate:"required,max=65536,min=1"`
	HTTPPort    int `validate:"required,max=65536,min=1"`
	JoinPort	int `validate:"max=65536,min=1"`
	DataDir     string `validate:"required,dir"`
	Bootstrap   bool
}
//required_with=JoinAddress,
type Config struct {
	RaftAddress net.Addr	//BindAddress + Raft Port
	HTTPAddress net.Addr	//BindAddress + HTTP Port
	JoinAddress string
	DataDir     string
	Bootstrap   bool
}
func ReadRawConfig() RawConfig {
	rawConf := cfg.ReadConf()

	return RawConfig{
		BindAddress: rawConf.RaftAddr,
		HTTPAddress: rawConf.HttpAddr,
		JoinAddress: rawConf.JoinAddr,
		RaftPort: rawConf.RaftPort,
		HTTPPort: rawConf.HttpPort,
		JoinPort: rawConf.JoinPort,
		DataDir: rawConf.RaftDataDir,
		Bootstrap: rawConf.Bootstrap,
	}
}
/**
 * This method validates the raw config and returns the final config
 */
func (rawConfig *RawConfig) ValidateConfig() (*Config, []validator.FieldError) {
	//Errors array in which all errors get saved
	errors := make([]validator.FieldError,0)
	//playground validator
	v := validator.New()
	err := v.Struct(rawConfig)
	//check for errors
	if err != nil {
		//loop through errors
		for _, fieldErr := range err.(validator.ValidationErrors) {
			if fieldErr.Tag() == "dir"{
				direrr := createDirectory(rawConfig.DataDir)
				if direrr == nil{
					continue
				}
			}
			errors = append(errors, fieldErr)
		}
		//return errors
		if errors != nil {
			return nil, errors
		}
	}
	//parse ip Address
	bindAddr := net.ParseIP(rawConfig.BindAddress)
	//create new tcpaddr from bindaddr and raftport
	raftAddr := &net.TCPAddr{
		IP: bindAddr,
		Port: rawConfig.RaftPort,
	}
	//create new tcpaddr from bindAddr and httpport
	httpAddress := net.ParseIP(rawConfig.HTTPAddress)
	httpAddr := &net.TCPAddr{
		IP: httpAddress,
		Port: rawConfig.HTTPPort,
	}

	//create new joinaddr from joinaddr and joinport
	joinAddr := ""
	if rawConfig.JoinAddress != "1" && rawConfig.JoinPort != 1  {
		joinAddr = rawConfig.JoinAddress + ":" + strconv.Itoa(rawConfig.JoinPort)
	}
	//create config
	config := &Config{
		RaftAddress: raftAddr,
		HTTPAddress: httpAddr,
		JoinAddress: joinAddr,
		DataDir:     rawConfig.DataDir,
		Bootstrap:   rawConfig.Bootstrap,
	}
	//return config
	return config, nil
}

func createDirectory(dir string) error{
	log.Print("creating directory "+dir)
	err := os.Mkdir(dir,0755)
	if err != nil{
		return err
	}
	return nil
}

/**
 * (unfinished) Validate config without go-playground/validator.v9
 */
/*func (rawConfig *RawConfig) validateConfig() (*Config, error) {

	var bindAddr net.IP
	bindAddr = net.ParseIP(rawConfig.BindAddress)
	if bindAddr == nil {
		log.Fatal("bind-Address could not be resolved")
	}

	var joinAddr net.IP
	joinAddr = net.ParseIP(rawConfig.JoinAddress)
	if joinAddr == nil {
		log.Fatal("join-Address could not be resolved")
	}

	//Check ports
	if rawConfig.RaftPort < 1 || rawConfig.RaftPort > 65536 {
		log.Fatal("Raft Port is invalid")
	}
	if rawConfig.HTTPPort < 1 || rawConfig.HTTPPort > 65536 {
		log.Fatal("HTTP Port is invalid")
	}

	//Construct addresses
	raftAddr := &net.TCPAddr{
		IP:   bindAddr,
		Port: rawConfig.RaftPort,
		Zone: "",
	}
	httpAddr := &net.TCPAddr{
		IP: bindAddr,
		Port: rawConfig.HTTPPort,
		Zone: "",
	}

	dataDir, err := filepath.Abs(rawConfig.DataDir)
	if err != nil {
		log.Fatal("data directory not valid")
	}

	config := &Config{
		RaftAddress: raftAddr,
		HTTPAddress: httpAddr,
		JoinAddress: string(joinAddr),
		DataDir:     dataDir,
		Bootstrap:   false,
	}
	return config
}*/