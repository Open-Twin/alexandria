package loadbalancing

import (
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var lbUrl = "http://127.0.0.1:8080/"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// Healthcheck Part

func TestHealthchecksSendPing(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"127.0.0.1": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := HealthCheck{
		Nodes:     &nodes,
		Interval:  10,
		CheckType: PingCheck,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(50)
	nodes = *hc.Nodes
	if nodes["127.0.0.1"].Healthy == false {
		t.Errorf("Sending ping healthcheck did not work: %v", nodes)
	}
}

func TestHealthchecksSendPingNodeOffline(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"12.12.12.12": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := HealthCheck{
		Nodes:     &nodes,
		Interval:  10,
		CheckType: PingCheck,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(50)
	nodes = *hc.Nodes
	if nodes["12.12.12.12"].Healthy == true {
		t.Errorf("Ping healthcheck falesly reported node as online: %v", nodes)
	}
}

func TestHealthchecksSendHttp(t *testing.T) {
	StartLoadReporting()

	nodes := map[storage.Ip]dns.NodeHealth{"127.0.0.1": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := HealthCheck{
		Nodes:     &nodes,
		Interval:  10,
		CheckType: HttpCheck,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(50)
	nodes = *hc.Nodes
	if nodes["127.0.0.1"].Healthy == false {
		t.Errorf("Sending http healthcheck did not work: %v", nodes)
	}
}

func TestHealthchecksSendHttpNodeOffline(t *testing.T) {
	nodes := map[storage.Ip]dns.NodeHealth{"12.12.12.12": {
		Healthy:     false,
		Connections: 0,
	}}

	hc := HealthCheck{
		Nodes:     &nodes,
		Interval:  10,
		CheckType: HttpCheck,
	}
	hc.ScheduleHealthChecks()

	time.Sleep(50)
	nodes = *hc.Nodes
	if nodes["127.0.0.1"].Healthy == true {
		t.Errorf("Http healthcheck falesly reported node as online: %v", nodes)
	}
}

// AlexandriaBalancer Part

func TestLoadbalancerSignupEndpoint(t *testing.T) {
	loadbalancer := AlexandriaBalancer{
		DnsPort: 53,
	}
	loadbalancer.StartAlexandriaLoadbalancer()

	data := url.Values{
		"ip": {"127.0.0.1"},
	}

	resp, err := http.PostForm(lbUrl+"register", data)
	if err != nil {
		t.Errorf("Fehler: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Fehler: %v", err)
	}

	if string(body) != "succesfully added" {
		t.Errorf("Adding node didn't work: %v", string(body))
	}
}

func TestLoadbalancerRequest(t *testing.T) {

}

func TestLoadbalancerServerGoesDown(t *testing.T) {

}

func TestLoadbalancerNoServerAdded(t *testing.T) {

}

func TestServerNoResponse(t *testing.T) {

}

/*
func startDnsServer() {
	// Auf Port 8333 hören
	connect, err := net.ListenPacket("udp", ":8333")

	if err != nil {
		fmt.Println(err)
	}
	defer connect.Close()

	for {
		msg := make([]byte, 512)
		// read message
		_, addr, err := connect.ReadFrom(msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Message from client: %s", addr.String())

		// Send answer in new thread
		go serve(connect, addr, []byte("this is the dns speaking. over"))
	}
}

func serve(conn net.PacketConn, addr net.Addr, msg []byte) {
	// send answer
	conn.WriteTo(msg, addr)
}

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

func TestStartServer(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("127.0.0.1:8333")
	answer := sendRequest("127.0.0.1:1212", t)
	fmt.Printf("Champagner: %s", answer)
	if !strings.HasPrefix(string(answer), "Message fowarded to: ") {
		t.Errorf("Wrong answer from dns-server: %s", string(answer))
	}
}

func sendRequest(ip string, t *testing.T) string {
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
/*
	receiverAddr, _ := net.ResolveUDPAddr("udp", ip)
	target, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	a, err := target.WriteToUDP([]byte("dejan.com"), receiverAddr)
	fmt.Printf("Holandese: %s", strconv.Itoa(a))

	if err != nil {
		t.Errorf("Bla: %s", err)
	}

	return strconv.Itoa(a)
	//return answer[0]
}

func TestResponse(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("127.0.0.1:8333")

	answer := string(sendRequest("127.0.0.1:1212", t))
	if answer != "this is the dns speaking. over" {
		t.Errorf("Wrong answer from dns server: %s", answer)
	}
}

func TestNoServerAdded(t *testing.T) {
	StartAlexandriaLoadbalancer(1212)

	answer := sendRequest("127.0.0.1:1212", t)
	fmt.Println(answer)
}

*/
