package storage

import "github.com/Open-Twin/alexandria/dns"
// Predefined variables for the usage in this class
//metadata
type service = string
type ip = string
type key = string
type value = string
//dns
type hostname = string
type record = dns.DNSResourceRecord

type Metadata struct {
	Dnsormetadata bool
	Service       string
	Ip            string
	Type          string
	Key           string
	Value         string
}
type Dnsresource struct {
	Dnsormetadata  bool
	Hostname       string
	Ip             string
	RequestType    string
	ResourceRecord dns.DNSResourceRecord
}

