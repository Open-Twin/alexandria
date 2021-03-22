package loadbalancing

func main() {
	loadbalancer := AlexandriaBalancer{
		DnsPort:             53,
		HealthCheckInterval: 30 * 1000,
	}
	loadbalancer.StartAlexandriaLoadbalancer()
}
