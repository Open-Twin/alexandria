package cfg

import (
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
	"testing"
)

func TestWelcome(t *testing.T) {
	fmt.Print("HALLO?!")
	conf := cfg.ReadConf()
	fmt.Println(conf.Welcome)
	if conf.Welcome != "Hello, I am a server." {
		t.Errorf("Welcome displays is: %s", conf.Welcome)
	}
}

func TestExampleTimeout(t *testing.T) {
	conf := cfg.ReadConf()
	if conf.ExampleTimeout != 10 {
		t.Errorf("ExampleTimeout is: %d", conf.ExampleTimeout)
	}
}

func TestEndpointPort(t *testing.T) {
	conf := cfg.ReadConf()
	if conf.EndpointPort != 1234 {
		t.Errorf("EndpointPort ist: %d", conf.EndpointPort)
	}
}
