package loadbalancing

import (
	"fmt"
	"github.com/Open-Twin/alexandria/communication"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type AlexandriaBalancer struct {
	DnsPort int
	nodes   map[storage.Ip]dns.NodeHealth
	pointer int
	lock    sync.RWMutex
}

func (lb *AlexandriaBalancer) StartAlexandriaLoadbalancer() {
	lb.nodes = make(map[storage.Ip]dns.NodeHealth)
	lb.pointer = 0
	lb.lock = sync.RWMutex{}

	lb.startSignupEndpoint()

	hc := HealthCheck{
		Nodes:     &lb.nodes,
		Interval:  5000,
		CheckType: HttpCheck,
	}
	hc.ScheduleHealthChecks()

	udpServer := communication.UDPServer{
		Address: []byte{0, 0, 0, 0},
		Port:    lb.DnsPort,
	}

	// Listen for connections
	go udpServer.Start(func(addr net.Addr, msg []byte) []byte {
		// Run the method for every message received
		go lb.forwardMsg(addr, msg)
		return []byte("request forwarded")
	})
}

func (balancer *AlexandriaBalancer) startSignupEndpoint() {
	http.HandleFunc("", balancer.addAlexandriaNode)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatalf("Signup Endpoint for Loadbalancer could not be started: %s", err.Error())
	}
}

func (balancer *AlexandriaBalancer) addAlexandriaNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ip := r.Form["ip"][0]

	balancer.lock.Lock()
	balancer.nodes[ip] = dns.NodeHealth{
		Healthy:     false,
		Connections: 0,
	}
	balancer.lock.Unlock()

	w.Write([]byte("succesfully added"))
}

/**
Removes a node from the loadbalanced dns nodes
*/
/*
func (balancer *AlexandriaBalancer) removeDns(dnsIp string) {
	balancer.lock.Lock()
	defer balancer.lock.Unlock()

	index := -1
	// search for item in list
	for i := 0; i < len(balancer.dnsservers); i++ {
		if balancer.dnsservers[i] == dnsIp {
			index = i
			i = len(balancer.dnsservers)
		}
	}

	// If the itmes was found
	if index != -1 {
		// append everthing before and after the item
		balancer.dnsservers = append(balancer.dnsservers[:index], balancer.dnsservers[index+1:]...)
	}
}*/

/**
Returns the next address in the list of loadbalanced nodes
*/
func (balancer *AlexandriaBalancer) nextAddr() string {
	balancer.lock.Lock()
	defer balancer.lock.Unlock()

	i := 0
	for ip, health := range balancer.nodes {
		if i == balancer.pointer && health.Healthy == true {
			return ip
		} else {
			i += 1
		}
		if i > len(balancer.nodes) {
			i = 0
			break
		}
	}

	balancer.pointer = i

	return ""
}

/**
Forwards an incoming message to a dns node
*/
func (balancer *AlexandriaBalancer) forwardMsg(source net.Addr, msg []byte) {
	fmt.Println("Message received: " + string(msg))

	if len(balancer.nodes) == 0 {
		log.Println("No dns nodes to forward to.")
		return
	}

	adrentik := balancer.nextAddr()
	if adrentik == "" {
		log.Println("No healthy nodes available.")
	}

	receiverAddr, err := net.ResolveUDPAddr("udp", adrentik+":"+strconv.Itoa(balancer.DnsPort))
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

	fmt.Printf("Message forwareded to: %s:%d\n", adrentik, balancer.DnsPort)
}
