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

func (hc *HealthCheck) ScheduleHealthChecks(interval int64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for range ticker.C {
			hc.loopNodes()
		}
	}()
}

func (hc *HealthCheck) loopNodes() {
	for i := range hc.nodes {
		n := &hc.nodes[i]
		n.sendCheck()
	}
}

func (node *Node) sendCheck() {
	resp, err := http.Get("http://" + node.ip + ":8080/load")
	if err != nil {
		node.healthy = false
		fmt.Printf("Node %s could not be reached: %s\n", node.ip, err)
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	node.healthy = true
	fmt.Printf("Load received %v %%\n", respBody)
}

func (hc *HealthCheck) AddNode(node string) {
	hc.nodes = append(hc.nodes, Node{node, false})
}

func (hc *HealthCheck) RemoveNode(node string) {
	index := -1
	// search for item in list
	for i, n := range hc.nodes {
		if n.ip == node {
			index = i
			break
		}
	}

	if index != -1 {
		hc.nodes = append(hc.nodes[:index], hc.nodes[index+1:]...)
	}
}

func (hc *HealthCheck) IsHealthy(node string) bool {
	for _, n := range hc.nodes {
		if n.ip == node {
			return n.healthy
		}
	}
	return false
}
