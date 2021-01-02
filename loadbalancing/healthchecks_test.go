package loadbalancing

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	startLoadReporting()
	code := m.Run()
	os.Exit(code)
}

func TestSendHealthchecks(t *testing.T) {
	hc := HealthCheck{[]Node{{"127.0.0.1", false}}}
	hc.ScheduleHealthChecks(10)
	time.Sleep(4 * time.Second)
	if hc.nodes[0].healthy == false {
		t.Errorf("Node not healthy: %s", hc.nodes[0].ip)
	}
}
