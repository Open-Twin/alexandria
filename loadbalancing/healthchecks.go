package loadbalancing

import (
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
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
	Nodes     map[storage.Ip]dns.NodeHealth
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
	log.Info().Msg("Running healthchecks")
	log.Debug().Msgf("Nodes to be healthchecked: %v", hc.Nodes)
	for ip := range hc.Nodes {
		node := hc.Nodes[ip]
		if hc.CheckType == HttpCheck {
			sendHttpCheck(ip, &node)
		} else if hc.CheckType == PingCheck {
			sendPingCheck(ip, &node)
		}
		hc.Nodes[ip] = node
	}
}

func sendHttpCheck(ip string, node *dns.NodeHealth) {
	resp, err := http.Get("http://" + ip + ":8080/health/")
	if err != nil {
		node.Healthy = false
		log.Info().Msgf("Node %s could not be reached: %s\n", ip, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		log.Info().Msgf("Node %s healthy. Response: %v\n", ip, string(respBody))
		node.Healthy = true
	} else {
		log.Info().Msgf("Node %s responded with bad status code %v: %s\n", ip, resp.StatusCode, resp.Status)
		node.Healthy = false
	}
}

func sendPingCheck(ip string, node *dns.NodeHealth) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Warn().Msgf("Error on creating pinger: %s\n", err.Error())
		return
	}

	os := runtime.GOOS
	if os == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = 3
	pinger.Interval = 10
	err = pinger.Run()
	if err != nil {
		log.Warn().Msgf("Error on sending ping: %s\n", err.Error())
		return
	}

	stats := pinger.Statistics()

	if stats.PacketsRecv > 1 {
		log.Info().Msgf("Node %s healthy. Statistics: %+v\n", ip, stats)
		node.Healthy = true
	} else {
		log.Info().Msgf("Node %s could not be reached. Statistics: %+v\n", ip, stats)
		node.Healthy = false
	}
}
