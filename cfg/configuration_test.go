package cfg

import (
	"os"
	"testing"
)

// Use following environment variable configuration for testing:
// HOSTNAME=peter;UDP_PORT=1234;TCP_PORT=4321;LOG_LEVEL=1;HTTP_ADDR=127.0.0.1;RAFT_ADDR=127.0.0.1;HTTP_PORT=11000;RAFT_PORT=12000;RAFT_DATA_DIR=/usr/bruh/awesomefolder;BOOTSTRAP=True

var conf Config

func TestMain(m *testing.M) {
	SetTestingConf()
	conf = ReadConf()
	code := m.Run()
	os.Exit(code)
}

func SetTestingConf() {
	os.Setenv("HOSTNAME", "peter")
	os.Setenv("UDP_PORT", "1234")
	os.Setenv("TCP_PORT", "4321")
	os.Setenv("LOG_LEVEL", "1")
	os.Setenv("HTTP_ADDR", "127.0.0.1")
	os.Setenv("RAFT_ADDR", "127.0.0.1")
	os.Setenv("HTTP_PORT", "11000")
	os.Setenv("RAFT_PORT", "12000")
	os.Setenv("JOIN_PORT", "13000")
	os.Setenv("RAFT_DATA_DIR", "/usr/bruh/awesomefolder")
	os.Setenv("BOOTSTRAP", "True")
}

func TestHostname(t *testing.T) {
	if conf.Hostname != "peter" {
		t.Errorf("%s is: %s", HOSTNAME, conf.Hostname)
	}
}

func TestHostnameEmpty(t *testing.T) {
	os.Unsetenv("Hostname")
	hostnameConf := ReadConf()

	if hostnameConf.Hostname != "" {
		t.Errorf("%s is: %s", HOSTNAME, conf.Hostname)
	}
}

func TestUdpPort(t *testing.T) {
	if conf.UdpPort != 1234 {
		t.Errorf("%s is: %d", UDP_PORT, conf.UdpPort)
	}
}

func TestTcpPort(t *testing.T) {
	if conf.TcpPort != 4321 {
		t.Errorf("%s is: %d", TCP_PORT, conf.TcpPort)
	}
}

func TestLoglevel(t *testing.T) {
	if conf.LogLevel != 1 {
		t.Errorf("%s is: %d", LOG_LEVEL, conf.LogLevel)
	}
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

func TestJoinPort(t *testing.T) {
	if conf.JoinPort != 13000 {
		t.Errorf("%s is: %d", JOIN_PORT, conf.RaftPort)
	}
}

func TestRaftDataDir(t *testing.T) {
	if conf.RaftDataDir != "/usr/bruh/awesomefolder" {
		t.Errorf("%s is: %s", RAFT_DATA_DIR, conf.RaftDataDir)
	}
}

func TestBootstrap(t *testing.T) {
	if conf.Bootstrap != true {
		t.Errorf("%s is: %t", BOOTSTRAP, conf.Bootstrap)
	}
}
