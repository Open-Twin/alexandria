package loadbalancing

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	startDnsServer()
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
		fmt.Println("Message from client: %s", addr)

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
}

func TestRemoveServer(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("192.168.0.1")
	lb.RemoveDns("192.168.0.1")
}

func TestRemoveNonExistent(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.RemoveDns("192.168.0.1")
}

func TestStartServer(t *testing.T) {
	StartAlexandriaLoadbalancer(1212)
	answer := sendRequest("127.0.0.1:1212", t)
	if !strings.HasPrefix(string(answer), "Message fowarded to: ") {
		t.Errorf("Wrong answer from dns-server: %s", string(answer))
	}
}

func TestResponse(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("127.0.0.1:83333")

	startDnsServer()
	answer := string(sendRequest("127.0.0.1:1212", t))
	if answer != "this is the dns speaking. over" {
		t.Errorf("Wrong answer from dns server: %s", answer)
	}
}

func TestNoServerAdded(t *testing.T) {
	StartAlexandriaLoadbalancer(1212)

	startDnsServer()
	answer := sendRequest("127.0.0.1:1212", t)
	fmt.Println(answer)
}

func sendRequest(ip string, t *testing.T) []byte {
	c, err := net.Dial("tcp", ip)
	if err != nil {
		t.Errorf("Error on connecting to loadbalancer: %s", err)
	}
	defer c.Close()

	_, err = fmt.Fprintf(c, "Hallo Loadbalancer")
	if err != nil {
		t.Errorf("Error on sending message: %s", err)
	}

	answer := make([]byte, 2048)
	_, err = bufio.NewReader(c).Read(answer)
	if err != nil {
		t.Errorf("Error on receiving message: %s", err)
	}

	return answer
}
