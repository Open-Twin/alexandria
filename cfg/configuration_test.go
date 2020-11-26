package cfg

import (
	"os"
	"testing"
)

var conf Config

func TestMain(m *testing.M) {
	conf = ReadConf()
	code := m.Run()
	os.Exit(code)
}

func TestHttpAddr(t *testing.T) {
	if conf.HttpAddr != "127.0.0.1" {
		t.Errorf("%s is: %s", HTTP_ADDR, conf.HttpAddr)
	}
}

func TestRaftAddr(t *testing.T) {
	if conf.RaftAddr != "127.0.0.1" {
		t.Errorf("%s is: %s", RAFT_ADDR, conf.RaftAddr)
	}
}

func TestHttpPort(t *testing.T) {
	if conf.HttpPort != 11000 {
		t.Errorf("%s is: %d", HTTP_PORT, conf.HttpPort)
	}
}

func TestRaftPort(t *testing.T) {
	if conf.RaftPort != 12000 {
		t.Errorf("%s is: %d", RAFT_PORT, conf.RaftPort)
	}
}

func TestRaftDataDir(t *testing.T) {
	if conf.RaftDataDir != "/usr/bruh/awesomefolder" {
		t.Errorf("%s is: %s", RAFT_DATA_DIR, conf.RaftDataDir)
	}
}

func TestBootstrap(t *testing.T) {
	if conf.Bootstrap != false {
		t.Errorf("%s is: %t", BOOTSTRAP, conf.Bootstrap)
	}
}
