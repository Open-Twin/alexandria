package loadbalancing

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"net"
)

//TODO: Comments and Multithreading Support

type AlexandriaBalancer struct {
	dnsservers []string
	pointer    int
	dnsport    int
}

func StartAlexandriaLoadbalancer(dnsport int) *AlexandriaBalancer {
	lb := AlexandriaBalancer{[]string{}, 0, dnsport}

	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    dnsport,
	}

	go udpServer.StartUDP(func(addr net.Addr, msg []byte) []byte {
		go lb.forwardMsg(msg)
		return []byte("dns forwarded")
	})

	return &lb
}

func (l *AlexandriaBalancer) AddDns(dnsIp string) {
	l.dnsservers = append(l.dnsservers, dnsIp)
}

func (l *AlexandriaBalancer) RemoveDns(dnsIp string) {
	index := -1
	// search for item in list
	for i := 0; i < len(l.dnsservers); i++ {
		if l.dnsservers[i] == dnsIp {
			index = i
			i = len(l.dnsservers)
		}
	}

	if index != -1 {
		l.dnsservers = append(l.dnsservers[:index], l.dnsservers[index+1:]...)
	}
}

func (l *AlexandriaBalancer) GetDnsEntries() []string {
	return l.dnsservers
}

func (l *AlexandriaBalancer) nextAddr() string {
	// Implementation of the loadbalancing
	l.pointer++
	if l.pointer >= len(l.dnsservers) {
		l.pointer = 0
	}

	address := l.dnsservers[l.pointer]
	return address
}

func (l *AlexandriaBalancer) forwardMsg(msg []byte) {
	if len(l.dnsservers)==0 {
		fmt.Println("No dns nodes to forward to.")
		return
	}

	adrentik := l.nextAddr()

	receiverAddr, err := net.ResolveUDPAddr("udp", adrentik + ":" + string(l.dnsport))
	if err != nil {
		fmt.Printf("Error on resolving client address : %s\n", err)
	}

	target, err := net.DialUDP("udp", nil, receiverAddr)
	if err != nil {
		fmt.Printf("Error on establishing client connection: %s\n", err)
	}

	_, err = target.Write(msg)
	if err != nil {
		fmt.Printf("Error on sending message to client: %s\n", err)
	}

	fmt.Printf("Message forwareded to: %s:%d\n", adrentik, l.dnsport)
}
