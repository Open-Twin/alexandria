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

func TestSendHealthchecks(t *testing.T) {
	hc := HealthCheck{}
	hc.AddNode("127.0.0.1")
	hc.ScheduleHealthChecks(500)
	time.Sleep(3 * time.Second)
	if hc.nodes[0].healthy == false {
		t.Errorf("Node not healthy: %s", hc.nodes[0].ip)
	}
}
