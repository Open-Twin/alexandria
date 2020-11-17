package loadbalancing

import (
	"fmt"
	"net"
	"testing"
)

func TestResponse(t *testing.T) {
	launchTestserver()
	loadbalanceAlexandriaNodes()
	if true {
		t.Errorf("Welcome displays is: %s", "lmao")
	}
}

func launchTestserver() {
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

		fmt.Println("Message from :", addr)

		// Send answer in new thread
		go serve(connect, addr, []byte("this is the dns speaking. over"))
	}
}

func serve(conn net.PacketConn, addr net.Addr, msg []byte) {
	// send answer
	conn.WriteTo(msg, addr)
}
