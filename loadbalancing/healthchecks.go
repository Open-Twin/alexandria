package loadbalancing

import (
	"fmt"
	"github.com/Open-Twin/alexandria/raft"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/go-ping/ping"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type NodeHealth struct {
	Healthy     bool
	Connections int
}

type HealthCheck struct {
	Node      raft.Node
	Interval  int
	CheckType CheckType
}

type CheckType int

const (
	HttpCheck CheckType = 0
	PingCheck CheckType = 1
)

func (hc *HealthCheck) ScheduleHealthChecks() {
	ticker := time.NewTicker(time.Duration(hc.Interval) * time.Millisecond)
	go func() {
		for range ticker.C {
			hc.loopNodes()
		}
	}()
}

func (hc *HealthCheck) loopNodes() {
	for ip := range hc.Node.Fsm.DnsRepo.LbInfo {
		node := hc.Node.Fsm.DnsRepo.LbInfo[ip]
		if hc.CheckType == HttpCheck {
			sendHttpCheck(ip, &node)
		} else if hc.CheckType == PingCheck {
			sendPingCheck(ip, &node)
		}
	}
}

func sendHttpCheck(ip string, node *NodeHealth) {
	resp, err := http.Get("http://" + ip + ":8080/health")
	if err != nil {
		node.Healthy = false
		fmt.Printf("Node %s could not be reached: %s\n", ip, err)
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("Node %s healthy. Response: %v\n", ip, string(respBody))
		node.Healthy = true
	} else {
		fmt.Printf("Node %s responded with bad status code %v: %s\n", ip, resp.StatusCode, resp.Status)
		node.Healthy = false
	}
}

func sendPingCheck(ip string, node *NodeHealth) {
	pinger, err := ping.NewPinger(ip)
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
		fmt.Printf("Node %s healthy. Statistics: %+v\n", ip, stats)
		node.Healthy = true
	} else {
		fmt.Printf("Node %s could not be reached. Statistics: %+v\n", ip, stats)
		node.Healthy = false
	}
}

func FindBestNode(hostname string, imsp *storage.StorageRepository) string {
	lowestConnections := 99999
	var lowestIp string
	for ip := range imsp.Entries[hostname] {
		if imsp.LbInfo[ip].Connections < lowestConnections {
			lowestConnections = imsp.LbInfo[ip].Connections
			lowestIp = ip
		}
	}

	nodeHealth := imsp.LbInfo[lowestIp]
	nodeHealth.Connections += 1
	imsp.LbInfo[lowestIp] = nodeHealth

	go func() {
		time.Sleep(time.Duration(imsp.Entries[hostname][lowestIp].TimeToLive) * time.Second)
		if imsp.Exists(hostname, lowestIp) {
			nodeHealth := imsp.LbInfo[lowestIp]
			nodeHealth.Connections -= 1
			imsp.LbInfo[lowestIp] = nodeHealth
		}
	}()

	return lowestIp
}
