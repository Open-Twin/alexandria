package loadbalancing

import (
	"fmt"
	"net"
	"testing"
)

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

func TestResponse(t *testing.T) {
	lb := StartAlexandriaLoadbalancer(1212)
	lb.AddDns("127.0.0.1")

	startDnsServer(t)
	sendRequest("127.0.0.1:1212", t)
}

func TestNoServerAdded(t *testing.T) {
	StartAlexandriaLoadbalancer(1212)

	startDnsServer(t)
	sendRequest("127.0.0.1:1212", t)
}

func sendRequest(ip string, t *testing.T) {
	c, err := net.Dial("tcp", ip)
	if err != nil {
		t.Errorf("Error on connecting to loadbalancer: %s", err)
	}
	_, err = fmt.Fprintf(c, "Hallo")
	if err != nil {
		t.Errorf("Error on sending message: %s", err)
	}
}

func startDnsServer(t *testing.T) {
	// Auf Port 8333 hören
	connect, err := net.ListenPacket("udp", ":1212")

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

		t.Logf("Message from client: %s", addr)

		// Send answer in new thread
		go serve(connect, addr, []byte("this is the dns speaking. over"))
	}
}

func serve(conn net.PacketConn, addr net.Addr, msg []byte) {
	// send answer
	conn.WriteTo(msg, addr)
}
