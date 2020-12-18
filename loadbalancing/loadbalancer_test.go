package loadbalancing

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go startDnsServer()
	code := m.Run()
	os.Exit(code)
}

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

func TestAddServer(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("192.168.0.1")
	if !equal(lb.GetDnsEntries(), []string{"192.168.0.1"}) {
		t.Errorf("Wrong entries after adding dns: %s", lb.GetDnsEntries())
	}
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

func TestRemoveServer(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("192.168.0.1")
	lb.RemoveDns("192.168.0.1")
	if !equal(lb.GetDnsEntries(), []string{}) {
		t.Errorf("Wrong entries after removing dns: %s", lb.GetDnsEntries())
	}
}

func TestRemoveNonExistent(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.RemoveDns("192.168.0.1")
}

func TestStartServer(t *testing.T) {
	StartAlexandriaLoadbalancer(1212)
	answer := sendRequest("172.0.0.1:1212", t)
	if !strings.HasPrefix(string(answer), "Message fowarded to: ") {
		t.Errorf("Wrong answer from dns-server: %s", string(answer))
	}
}

func sendRequest(ip string, t *testing.T) string {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", ip)
		},
	}
	answer, err := r.LookupHost(context.Background(), "www.example.com")
	if err!= nil {
		t.Errorf("Bla: %s", err)
	}
	fmt.Println(answer)

	return answer[0]
}

func TestResponse(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("127.0.0.1:83333")

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
