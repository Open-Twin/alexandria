package loadbalancing

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"net"
	"strconv"
	"sync"
)

type AlexandriaBalancer struct {
	dnsport    int
	dnsservers []string
	pointer    int
	lock       sync.RWMutex
}

/**
Starts listening for connections on the specified dns port
*/
func StartAlexandriaLoadbalancer(dnsport int) *AlexandriaBalancer {
	lb := AlexandriaBalancer{dnsport, []string{}, 0, sync.RWMutex{}}

	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    dnsport,
	}

	// Listen for connections
	go udpServer.StartUDP(func(addr net.Addr, msg []byte) []byte {
		// Run the method for every message received
		go lb.forwardMsg(addr, msg)
		return []byte("request forwarded")
	})

	return &lb
}

/**
Adds a node to the list of loadbalaced dns nodes
*/
func (l *AlexandriaBalancer) AddDns(dnsIp string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// append new node to list
	l.dnsservers = append(l.dnsservers, dnsIp)
}

/**
Removes a node from the loadbalanced dns nodes
 */
func (l *AlexandriaBalancer) RemoveDns(dnsIp string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	index := -1
	// search for item in list
	for i := 0; i < len(l.dnsservers); i++ {
		if l.dnsservers[i] == dnsIp {
			index = i
			i = len(l.dnsservers)
		}
	}

	// If the itmes was found
	if index != -1 {
		// append everthing before and after the item
		l.dnsservers = append(l.dnsservers[:index], l.dnsservers[index+1:]...)
	}
}

/**
Returns all loadbalanced dns entries
 */
func (l *AlexandriaBalancer) GetDnsEntries() []string {
	return l.dnsservers
}

/**
Returns the next address in the list of loadbalanced nodes
 */
func (l *AlexandriaBalancer) nextAddr() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	// implementation of the loadbalancing algorithm (round robin)
	// move the pointer one ahead
	l.pointer++
	// if the pointer is larger than the number of nodes it has to be reset
	if l.pointer >= len(l.dnsservers) {
		l.pointer = 0
	}

	address := l.dnsservers[l.pointer]
	return address
}

/**
Forwards an incoming message to a dns node
 */
func (l *AlexandriaBalancer) forwardMsg(source net.Addr, msg []byte) {
	fmt.Println("Message received: " + string(msg))

	if len(l.dnsservers) == 0 {
		fmt.Println("No dns nodes to forward to.")
		return
	}

	adrentik := l.nextAddr()

	receiverAddr, err := net.ResolveUDPAddr("udp", adrentik+":"+strconv.Itoa(l.dnsport))
	if err != nil {
		fmt.Printf("Error on resolving dns address : %s\n", err)
	}

	sourceAddr, err := net.ResolveUDPAddr("udp", source.String())
	if err != nil {
		fmt.Printf("Error on resolving client address : %s\n", err)
	}

	target, err := net.DialUDP("udp", sourceAddr, receiverAddr)
	if err != nil {
		fmt.Printf("Error on establishing dns connection: %s\n", err)
	}

	_, err = target.WriteToUDP(msg, receiverAddr)
	if err != nil {
		fmt.Printf("Error on sending message to dns: %s\n", err)
	}

	fmt.Printf("Message forwareded to: %s:%d\n", adrentik, l.dnsport)
}
