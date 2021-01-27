package raft_test

import (
	"fmt"
	"github.com/Open-Twin/alexandria/raft"
	"testing"
)
/*
Tests config validation for general errors if params are technically correct
 */
func TestValidateValidConfig(t *testing.T){
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.1",
		JoinAddress: "192.168.0.1",
		RaftPort:    1000,
		HTTPPort:    8080,
		DataDir:     "../raft",
		Bootstrap:   false,
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
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.f",
		JoinAddress: "192.168.0.1",
		RaftPort:    1000,
		HTTPPort:    8080,
		DataDir:     "../raft",
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
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.1",
		JoinAddress: "192.a.0.1",
		RaftPort:    1000,
		HTTPPort:    8080,
		DataDir:     "./raft",
		Bootstrap:   false,
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
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.1",
		JoinAddress: "192.168.0.1",
		RaftPort:    100000,
		HTTPPort:    8080,
		DataDir:     "./raft",
		Bootstrap:   false,
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
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.1",
		JoinAddress: "192.168.0.1",
		RaftPort:    1000,
		HTTPPort:    -8080,
		DataDir:     "./raft",
		Bootstrap:   false,
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
func TestValidateConfigUnvalidDataDir(t *testing.T){
	cfg := &raft.RawConfig{
		BindAddress: "192.168.0.1",
		JoinAddress: "192.168.0.1",
		RaftPort:    1000,
		HTTPPort:    -8080,
		DataDir:     "./raft/",
		Bootstrap:   false,
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