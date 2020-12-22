package loadbalancing

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	startLoadReporting()
	code := m.Run()
	os.Exit(code)
}

func TestStartHealthchecks(t *testing.T) {
	scheduleHealthChecks(10, []string{"127.0.0.1"})
}
