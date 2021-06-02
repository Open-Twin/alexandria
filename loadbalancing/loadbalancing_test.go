package loadbalancing_test

import (
	"context"
	"fmt"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Level(3))
	code := m.Run()
	os.Exit(code)
}

// HealthCheck Part

func TestAutomaticDeletion(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"12.12.12.12": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := loadbalancing.HealthCheck{
		Nodes:          nodes,
		Interval:       50 * time.Millisecond,
		CheckType:      loadbalancing.PingCheck,
		RemoveTimeout:  1 * time.Second,
		RequestTimeout: 20 * time.Millisecond,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(2 * time.Second)
	nodes = hc.Nodes
	if _, ok := nodes["12.12.12.12"]; ok {
		t.Error("Nodes was not deleted from map after timeout.")
	}
}

func TestHealthchecksSendPing(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"127.0.0.1": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := loadbalancing.HealthCheck{
		Nodes:          nodes,
		Interval:       30 * time.Millisecond,
		CheckType:      loadbalancing.PingCheck,
		RemoveTimeout:  5 * time.Second,
		RequestTimeout: 20 * time.Millisecond,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(400 * time.Millisecond)
	nodes = hc.Nodes
	if nodes["127.0.0.1"].Healthy == false {
		t.Errorf("Sending ping healthcheck did not work: %v", nodes)
	}
}

func TestHealthchecksSendPingNodeOffline(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"12.12.12.12": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := loadbalancing.HealthCheck{
		Nodes:          nodes,
		Interval:       30 * time.Millisecond,
		CheckType:      loadbalancing.PingCheck,
		RemoveTimeout:  2 * time.Second,
		RequestTimeout: 20 * time.Millisecond,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(400 * time.Millisecond)
	nodes = hc.Nodes
	if nodes["12.12.12.12"].Healthy == true {
		t.Errorf("Ping healthcheck falesly reported node as online: %v", nodes)
	}
}

func TestHealthchecksSendHttp(t *testing.T) {
	loadbalancing.StartLoadReporting("127.0.0.1", 8001)

	nodes := map[storage.Ip]dns.NodeHealth{"127.0.0.1": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := loadbalancing.HealthCheck{
		Nodes:          nodes,
		Interval:       30 * time.Millisecond,
		CheckType:      loadbalancing.HttpCheck,
		HttpPingPort:   8001,
		RemoveTimeout:  1 * time.Second,
		RequestTimeout: 20 * time.Millisecond,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(300 * time.Millisecond)
	nodes = hc.Nodes
	if nodes["127.0.0.1"].Healthy == false {
		t.Errorf("Sending http healthcheck did not work: %v", nodes)
	}
}

func TestHealthchecksSendHttpNodeOffline(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"12.12.12.12": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := loadbalancing.HealthCheck{
		Nodes:          nodes,
		Interval:       10 * time.Millisecond,
		CheckType:      loadbalancing.HttpCheck,
		HttpPingPort:   8080,
		RemoveTimeout:  1 * time.Second,
		RequestTimeout: 20 * time.Millisecond,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(50 * time.Millisecond)
	nodes = hc.Nodes
	if nodes["12.12.12.12"].Healthy == true {
		t.Errorf("Http healthcheck falesly reported node as online: %v", nodes)
	}
}

// AlexandriaBalancer Part
var dnsAnswer = "My name is dns."

func TestLoadbalancerNoServerAdded(t *testing.T) {
	loadbalancer := loadbalancing.AlexandriaBalancer{
		RegistrationPort:    14000,
		DnsPort:             10000,
		DnsApiPort:          14001,
		MetdataApiPort:      14002,
		HttpPingPort:        14003,
		HealthCheckInterval: 5 * time.Second,
		RemoveNodeTimeout:   5 * time.Second,
		RequestTimeout:      100 * time.Millisecond,
	}
	go loadbalancer.StartAlexandriaLoadbalancer()

	/*
		_, err := sendRequest("127.0.0.1:"+strconv.Itoa(10010), "www.dejan.com")
		if err == nil {
			fmt.Errorf("Crap")
		}
	*/
}

func TestLoadbalancerSignupEndpoint(t *testing.T) {
	lbUrl := "http://127.0.0.1:" + strconv.Itoa(14000)

	t.Logf("Signing in at loadbalancer")
	err := signinLocalhost(t, lbUrl)
	if err != nil {
		t.Errorf("Registration did not work: %v", err.Error())
	}
}

func TestLoadbalancerForwardRequest(t *testing.T) {
	dnsIp := "10.6.0.3"

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, dnsIp)
		},
	}
	ans, err := r.LookupIPAddr(context.Background(), "www.dejan.at")
	if err != nil {
		t.Errorf("Error on sending request: %v", err)
	}
	gesucht := "185.21.102.144"
	if ans[0].String() != gesucht {
		t.Errorf("Wrong ip returned: %v", ans[0])
	}
}

/*func TestLoadbalancerServerGoesDown(t *testing.T) {

}*/

func signinLocalhost(t *testing.T, lbUrl string) (err error) {
	resp, err := http.Get(lbUrl + "/signup")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "succesfully added" {
		return err
	}
	return nil
}

func startTestingDns(dnsPort int) error {
	connect, err := net.ListenPacket("udp", ":"+strconv.Itoa(dnsPort))

	if err != nil {
		return err
	}
	defer connect.Close()

	for {
		msg := make([]byte, 512)
		// read message
		_, addr, err := connect.ReadFrom(msg)
		if err != nil {
			return err
		}
		fmt.Printf("Message from client: %s", addr.String())

		// Send answer in new thread
		go func(conn net.PacketConn, addr net.Addr, msg string) {
			conn.WriteTo([]byte(msg), addr)
		}(connect, addr, dnsAnswer)
	}
}

func sendRequest(lbip string, request string) (string, error) {
	/*r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", ip)
		},
	}
	answer, err := r.LookupHost(context.Background(), "www.example.com")
	*/
	receiverAddr, _ := net.ResolveUDPAddr("udp", lbip)
	target, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	a, err := target.WriteToUDP([]byte(request), receiverAddr)
	fmt.Printf("Holandese: %s", strconv.Itoa(a))

	if err != nil {
		return "", err
	}

	return strconv.Itoa(a), nil
	//return answer[0] }
}

/*
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
*/
