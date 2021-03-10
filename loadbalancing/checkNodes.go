package loadbalancing

import (
	"fmt"
	"github.com/go-ping/ping"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type HealthCheck struct {
	nodes []Node
}

type Node struct {
	ip      string
	healthy bool
}

type CheckType int

const (
	HttpCheck CheckType = 0
	PingCheck CheckType = 1
)

func (hc *HealthCheck) ScheduleHealthChecks(interval int, checkType CheckType) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for range ticker.C {
			hc.loopNodes(checkType)
		}
	}()
}

func (hc *HealthCheck) loopNodes(checkType CheckType) {
	for i := range hc.nodes {
		n := &hc.nodes[i]
		if checkType == HttpCheck {
			n.sendHttpCheck()
		} else if checkType == PingCheck {
			n.sendPingCheck()
		}
	}
}

func (node *Node) sendHttpCheck() {
	resp, err := http.Get("http://" + node.ip + ":8080/health")
	if err != nil {
		node.healthy = false
		fmt.Printf("Node %s could not be reached: %s\n", node.ip, err)
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("Node %s healthy. Response: %v\n", node.ip, string(respBody))
		node.healthy = true
	} else {
		fmt.Printf("Node %s responded with bad status code %v: %s\n", node.ip, resp.StatusCode, resp.Status)
		node.healthy = false
	}
}

func (node *Node) sendPingCheck() {
	pinger, err := ping.NewPinger(node.ip)
	if err != nil {
		fmt.Printf("Error on creating pinger: %s\n", err.Error())
	}

	os := runtime.GOOS
	if os == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Error on sending ping: %s\n", err.Error())
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 1 {
		fmt.Printf("Node %s healthy. Statistics: %+v\n", node.ip, stats)
		node.healthy = true
	} else {
		fmt.Printf("Node %s could not be reached. Statistics: %+v\n", node.ip, stats)
		node.healthy = false
	}
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
