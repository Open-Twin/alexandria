package loadbalancing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HealthCheck struct {
	nodes []string
}

func scheduleHealthChecks(interval int64, nodes []string) *HealthCheck {
	hc := HealthCheck{
		nodes,
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for range ticker.C {
			loopNodes(hc.nodes)
		}
	}()
	return &hc
}

func loopNodes(nodes []string) {
	for _, s := range nodes {
		sendCheck(s)
	}
}

func sendCheck(ip string) {
	resp, err := http.Get(ip + "/load")
	if err != nil {
		fmt.Printf("Node %s could not be reached: %s", ip, err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Load received %v %%", respBody)
}

func (hc *HealthCheck) addNode(node string) {
	hc.nodes = append(hc.nodes, node)
}

func (hc *HealthCheck) removeNode(node string) {
	// TODO
}
