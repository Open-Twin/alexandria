package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	connect, err := net.ListenPacket("udp", ":53")

	if err != nil {
		fmt.Println(err)
	}
	defer connect.Close()

	fmt.Println("Listening...")
	for {
		msg := make([]byte, 512)
		n, addr, err := connect.ReadFrom(msg)
		if err!= nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("Message from %v: ", addr)

		received := string(msg[:])
		received = strings.TrimSpace(received)
		splitMessage := strings.Split(received, "\n")
		for _, data := range splitMessage {
			fmt.Print(data)
		}
		fmt.Print("\n")

		// Antwort senden in neuem Thread
		go serve(connect, addr, msg[:n])
	}
}

func serve(conn net.PacketConn, addr net.Addr, msg []byte) {
	msg = []byte("This is the dns speaking.")
	conn.WriteTo(msg, addr)
}
