package communication

import (
	"fmt"
	"net"
	"time"
)

func scheduleHealthChecks(dnsPort int, interval int64) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for range ticker.C {
			loopNodes(dnsPort)
		}
	}()

	// Stop after 10 seconds
	//time.Sleep(10 *time.Second)
	//ticker.Stop()
}

func loopNodes(dnsPort int) {
	nodes := []string{"12.12.12.12", "13.13.13.13."}

	for _, s := range nodes {
		sendCheck(s)
	}
}

func sendCheck(ip string, dnsPort int) {
	dnsConn, err := net.Dial("udp", ip+":"+string(dnsPort))
	if err != nil {
		fmt.Printf("We had an error.")
	}
	fmt.Fprintf(dnsConn, "are you ok?")
	dnsConn.Close()
	fmt.Print("Message forwareded to :", ip+":"+string(dnsPort))
}
