package cfg

import (
	"fmt"
	"log"
	"os"
	"testing"
)
var wdpath string
func TestMain(m *testing.M) {
	path, _ := os.Getwd()
	wdpath = path
	code := m.Run()
	os.Exit(code)
}

func SetTestingConf() {
	os.Setenv("HOSTNAME", "peter")
	os.Setenv("LOG_LEVEL", "1")
	os.Setenv("DATA_DIR", wdpath)
	os.Setenv("BOOTSTRAP", "True")
	os.Setenv("AUTO_JOIN", "True")
	os.Setenv("HEALTHCHECK_INTERVAL", "3000")
	os.Setenv("HTTP_ADDR", "127.0.0.1")
	os.Setenv("RAFT_ADDR", "127.0.0.1")
	os.Setenv("JOIN_ADDR", "127.0.0.1")
	os.Setenv("HTTP_PORT", "1234")
	os.Setenv("RAFT_PORT", "4321")
}

func TestConfSetCorrectly(t *testing.T) {
	SetTestingConf()
	cfg := ReadConf()
	if cfg.Hostname != "peter" {
		t.Errorf("Reading Hostname failed: %s\n", cfg.Hostname)
	}
	if cfg.LogLevel != 1 {
		t.Errorf("Reading LogLevel failed: %v\n", cfg.LogLevel)
	}
	if cfg.DataDir != wdpath {
		t.Errorf("Reading DataDir failed: %s\n", cfg.DataDir)
	}
	if cfg.Bootstrap != true {
		t.Errorf("Reading Bootstrap failed: %v\n", cfg.Bootstrap)
	}
	if cfg.Autojoin != true {
		t.Errorf("Reading AutoJoin failed: %v\n", cfg.Autojoin)
	}
	if cfg.HealthcheckInterval != 3000 {
		t.Errorf("Reading HealthcheckInterval failed: %v\n", cfg.HealthcheckInterval)
	}
	if cfg.HttpAddr.String() != "127.0.0.1:1234" {
		t.Errorf("Reading HttpAddress failed: %v\n", cfg.HttpAddr)
	}
	if cfg.RaftAddr.String() != "127.0.0.1:4321" {
		t.Errorf("Reading RaftAddress failed: %v\n", cfg.RaftAddr)
	}
	if cfg.JoinAddr.String() != "127.0.0.1:1234" {
		t.Errorf("Reading JoinAddress failed: %s\n", cfg.JoinAddr)
	}
}

/*
Tests config validation for general errors if params are technically correct
*/
func TestValidateValidConfig(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		MetaApiAddr: "1.2.3.4",
		DnsApiAddr: "1.2.3.4",
		DnsAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: 8000,
		MetaApiPort: 20000,
		DnsApiPort: 10000,
		UdpPort: 9000,
		DnsPort: 53,
	}
	_, err := validateConfig(cfg)
	if err != nil {
		t.Errorf("cfg.validateConfig failed, expected %v but got %v", nil, err)
	}
}

/*
Tests config validation for errors at 'BindAddress'
*/
func TestValidateConfigUnvalidBindAddress(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.f",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "RaftAddr"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			log.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound{
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}
/*
Tests config for errors at 'JoinAddress'
*/
func TestValidateConfigUnvalidJoinAddress(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.f.4",
		RaftPort: 7000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "JoinAddr"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}
/*
Tests config validation for errors at 'RaftPort'
*/
func TestValidateConfigUnvalidRaftPort(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 100000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "RaftPort"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}
/*
Tests config validation for errors at 'HTTPPort'
*/
func TestValidateConfigUnvalidHTTPPort(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: -2,
	}
	errorFound := false
	expectedResult := "HttpPort"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}

func TestValidateConfigInvalidHealthcheckInterval(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: 1,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: -20,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "HealthcheckInterval"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}

func TestValidateConfigInvalidLoglevel(t *testing.T){
	cfg := rawConfig{
		Hostname: "adincarik",
		LogLevel: -5,
		DataDir: wdpath,
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "LogLevel"
	result, errors := validateConfig(cfg)
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}

/*
Tests config validation for errors at 'DataDir'
*/
/*func TestValidateConfigUnvalidDataDir(t *testing.T){
	cfg := rawConfig{
		Hostname: "adin carik",
		LogLevel: 1,
		DataDir: "",
		Bootstrap: false,
		Autojoin: false,
		HealthcheckInterval: 2000,
		RaftAddr: "1.2.3.4",
		HttpAddr: "1.2.3.4",
		JoinAddr: "1.2.3.4",
		RaftPort: 7000,
		HttpPort: 8000,
	}
	errorFound := false
	expectedResult := "DataDir"
	result, errors := cfg.ValidateConfig()
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
			if err.Field()==expectedResult{
				errorFound = true
			}
		}
	}
	if !errorFound {
		t.Errorf("cfg.validateConfig , expected an error at "+expectedResult+" but got %v", result)
	}
}
*/

