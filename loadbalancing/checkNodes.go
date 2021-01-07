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
	for i, _ := range hc.nodes {
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
	// TODO
}
