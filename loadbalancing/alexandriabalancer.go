package loadbalancing

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"net"
)

type AlexandriaBalancer struct {
	dnsservers []string
	pointer    int
	dnsport    int
}

func StartAlexandriaLoadbalancer(dnsport int) AlexandriaBalancer {
	lb := AlexandriaBalancer{[]string{}, 0, dnsport}

	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    dnsport,
	}

	go udpServer.StartUDP(func(addr net.Addr, msg []byte) []byte {
		go lb.forwardMsg(msg)
		return []byte("dns forwarded")
	})

	return lb
}

func (l AlexandriaBalancer) AddDns(dnsIp string) {
	l.dnsservers = append(l.dnsservers, dnsIp)
}

func (l AlexandriaBalancer) RemoveDns(dnsIp string) {
	index := -1
	// search for item in list
	for i := 0; i < len(l.dnsservers); i++ {
		if l.dnsservers[i] == dnsIp {
			index = i
			i = len(l.dnsservers)
		}
	}

	if index != -1 {
		l.dnsservers = append(l.dnsservers[:index], l.dnsservers[index:]...)
	}
}

func (l AlexandriaBalancer) nextAddr() string {
	// Implementation of the loadbalancing
	l.pointer++
	if l.pointer > len(l.dnsservers) {
		l.pointer = 0
	}

	address := l.dnsservers[l.pointer]
	//adrentik := net.IPAddr{IP: net.ParseIP(address + string(dnsport))}
	//return &adrentik
	return address
}

func (l AlexandriaBalancer) forwardMsg(msg []byte) {
	adrentik := l.nextAddr()

	receiverAddr, err := net.ResolveUDPAddr("udp", adrentik)
	if err != nil {
		fmt.Printf("Problem")
	}

	target, err := net.DialUDP("udp", nil, receiverAddr)
	if err != nil {
		fmt.Printf("Problem")
	}

	_, err = target.WriteToUDP(msg, receiverAddr)
	if err != nil {
		fmt.Printf("We had an error.")
	}

	fmt.Print("Message forwareded to :", adrentik)
}
