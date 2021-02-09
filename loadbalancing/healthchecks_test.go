package loadbalancing

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	StartLoadReporting()
	code := m.Run()
	os.Exit(code)
}

func TestSendHttpHealthchecks(t *testing.T) {
	hc := HealthCheck{}
	hc.AddNode("127.0.0.1")
	hc.ScheduleHealthChecks(500, HttpCheck)
	time.Sleep(3 * time.Second)
	if hc.nodes[0].healthy == false {
		t.Errorf("Node not healthy: %s", hc.nodes[0].ip)
	}
}

func TestAddServer(t *testing.T) {
	hc := HealthCheck{}
	hc.AddNode("127.0.0.1")
	if !equal(hc.nodes, []Node{{"127.0.0.1", false}}) {
		t.Errorf("Wrong entries after adding dns: %v", hc.nodes)
	}
}

func equal(a, b []Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestRemoveOneServerInList(t *testing.T) {
	hc := HealthCheck{}
	hc.AddNode("192.168.0.1")
	hc.RemoveNode("192.168.0.1")
	if !equal(hc.nodes, []Node{}) {
		t.Errorf("Wrong entries after removing dns: %v", hc.nodes)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	hc := HealthCheck{}
	hc.RemoveNode("192.168.0.1")
}
