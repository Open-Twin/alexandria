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

type HealthCheck struct {
	Nodes          map[storage.Ip]dns.NodeHealth
	Interval       time.Duration
	CheckType      CheckType
	RequestTimeout time.Duration
	RemoveTimeout  time.Duration
}

type CheckType int

const (
	HttpCheck CheckType = 0
	PingCheck CheckType = 1
)

func (hc *HealthCheck) ScheduleHealthChecks() {
	if hc.RequestTimeout == 0 {
		hc.RequestTimeout = 1 * time.Second
	}

	ticker := time.NewTicker(hc.Interval)
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
			sendHttpCheck(ip, &node, hc.RequestTimeout)
		} else if hc.CheckType == PingCheck {
			sendPingCheck(ip, &node, hc.RequestTimeout)
		}
		hc.Nodes[ip] = node

		if node.Healthy == false {
			hc.checkForDeletion(ip)
		}
	}
}

func sendHttpCheck(ip string, node *dns.NodeHealth, timeout time.Duration) {
	client := http.Client{
		Timeout: timeout,
	}
	//TODO: port?
	resp, err := client.Get("http://" + ip + ":8080/health")
	if err != nil {
		node.Healthy = false
		log.Info().Msgf("Node %s could not be reached: %s\n", ip, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		log.Info().Msgf("Node %s healthy. Response: %v\n", ip, string(respBody))
		node.Healthy = true
		node.LastOnline = time.Now()
	} else {
		log.Info().Msgf("Node %s responded with bad status code %v: %s\n", ip, resp.StatusCode, resp.Status)
		node.Healthy = false
	}
}

func sendPingCheck(ip string, node *dns.NodeHealth, timeout time.Duration) {
	log.Info().Msg("IPIDIOT: "+ip)
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
	pinger.Timeout = timeout

	err = pinger.Run()

	if err != nil {
		log.Warn().Msgf("Error on sending ping: %s\n", err.Error())
		return
	}

	stats := pinger.Statistics()

	if stats.PacketsRecv > 1 {
		log.Info().Msgf("Node %s healthy. Statistics: %+v\n", ip, stats)
		node.Healthy = true
		node.LastOnline = time.Now()
	} else {
		log.Info().Msgf("Node %s could not be reached. Statistics: %+v\n", ip, stats)
		node.Healthy = false

	}
}

func (hc *HealthCheck) checkForDeletion(ip storage.Ip) {
	node := hc.Nodes[ip]

	if node.LastOnline.IsZero() {
		node.LastOnline = time.Now()
		hc.Nodes[ip] = node
	} else if time.Now().After(node.LastOnline.Add(hc.RemoveTimeout)) {
		delete(hc.Nodes, ip)
	}
}
