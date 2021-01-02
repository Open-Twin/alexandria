package loadbalancing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HealthCheck struct {
	nodes []Node
}

type Node struct {
	ip      string
	healthy bool
}

func (hc HealthCheck) ScheduleHealthChecks(interval int64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for range ticker.C {
			hc.loopNodes()
		}
	}()
}

func (hc HealthCheck) loopNodes() {
	for _, node := range hc.nodes {
		node.sendCheck()
	}
}

func (node Node) sendCheck() {
	resp, err := http.Get("http://" + node.ip + ":8080/load")
	if err != nil {
		node.healthy = false
		fmt.Printf("Node %s could not be reached: %s", node.ip, err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	node.healthy = true
	fmt.Printf("Load received %v %%", respBody)
}

func (hc *HealthCheck) AddNode(node Node) {
	hc.nodes = append(hc.nodes, node)
}

func (hc *HealthCheck) RemoveNode(node Node) {
	// TODO
}
