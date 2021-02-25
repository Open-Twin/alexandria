package config_test

import (
	"fmt"
	"github.com/Open-Twin/alexandria/raft/config"
	"testing"
)
/*
Tests config validation for general errors if params are technically correct
 */
func TestValidateValidConfig(t *testing.T){
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.4",
		JoinAddress: "1.2.3.4",
		RaftPort:    1000,
		HTTPPort:    8080,
		JoinPort:    8000,
		DataDir:     "../config",
		Bootstrap:   false,
		AutoJoin: false,
		AutojoinPort: 9000,
	}
	_, err := cfg.ValidateConfig()
	if err != nil {
		t.Errorf("cfg.validateConfig failed, expected %v but got %v", nil, err)
	}
}

/*
Tests config validation for errors at 'BindAddress'
 */
func TestValidateConfigUnvalidBindAddress(t *testing.T){
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.f",
		JoinAddress: "1.2.3.4",
		RaftPort:    1000,
		HTTPPort:    8080,
		JoinPort:	 8000,
		DataDir:     "../config",
		Bootstrap:   false,
	}
	errorFound := false
	expectedResult := "BindAddress"
	result, errors := cfg.ValidateConfig()
	if errors != nil {
		for _, err := range errors {
			fmt.Println("Errors: "+err.Field())
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
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.4",
		JoinAddress: "1.a.3.4",
		RaftPort:    1000,
		HTTPPort:    8080,
		JoinPort:	 8000,
		DataDir:     "../config",
		Bootstrap:   false,
		AutoJoin: false,
		AutojoinPort: 9000,
	}
	errorFound := false
	expectedResult := "JoinAddress"
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
/*
Tests config validation for errors at 'RaftPort'
 */
func TestValidateConfigUnvalidRaftPort(t *testing.T){
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.4",
		JoinAddress: "1.2.3.4",
		RaftPort:    100000,
		HTTPPort:    8080,
		JoinPort:	 8000,
		DataDir:     "../config",
		Bootstrap:   false,
		AutoJoin: false,
		AutojoinPort: 9000,
	}
	errorFound := false
	expectedResult := "RaftPort"
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
/*
Tests config validation for errors at 'HTTPPort'
 */
func TestValidateConfigUnvalidHTTPPort(t *testing.T){
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.4",
		JoinAddress: "1.2.3.4",
		RaftPort:    1000,
		HTTPPort:    -8080,
		JoinPort:	 8000,
		DataDir:     "../config",
		Bootstrap:   false,
		AutoJoin: false,
		AutojoinPort: 9000,
	}
	errorFound := false
	expectedResult := "HTTPPort"
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
/*
Tests config validation for errors at 'DataDir'
 */
/*func TestValidateConfigUnvalidDataDir(t *testing.T){
	cfg := &config.RawConfig{
		BindAddress: "1.2.3.4",
		JoinAddress: "1.2.3.4",
		RaftPort:    1000,
		HTTPPort:    -8080,
		JoinPort:	 8000,
		DataDir:     "./unvalid",
		Bootstrap:   false,
		AutoJoin: false,
		AutojoinPort: 9000,
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
}*/