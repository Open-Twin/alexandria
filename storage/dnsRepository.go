package storage

import (
	"errors"
	"github.com/Open-Twin/alexandria/dns"
	"sync"
	"time"
)

/*
	This class is also seen as the "Repository" it contains the basic CRUD (= Create, Read, Update, Delete) functionalities
	for the Metadata storage, which is nearly linked to the API.
*/

// Creating the MetadataRepository interface which predefines the needed methods for the storage
type DNSRepository interface {
	Exists(service, ip, key string) bool
	Create(service, ip, key, value string) error
	Read(service, ip, key string) (string, error)
	Update(service, ip, key, value string) error
	Delete(service, ip, key string) error
}

// Creating a structure for the Metadata containing the necessary variables
type StorageRepository struct {
	Entries map[hostname]map[Ip]record
	// Creating a mutex onto the Metadata variable in order to handle threads
	mutex  sync.RWMutex
	LbInfo map[Ip]dns.NodeHealth
}

func NewInMemoryDNSStorageRepository() *StorageRepository {
	entries := make(map[hostname]map[Ip]record)
	info := make(map[Ip]dns.NodeHealth)
	return &StorageRepository{
		Entries: entries,
		LbInfo:  info,
	}
}

// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *StorageRepository) Exists(hostname, ip string) bool {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	_, ok := imsp.Entries[hostname][ip]
	return ok
}

// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *StorageRepository) ExistsHostname(hostname string) bool {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	_, ok := imsp.Entries[hostname]
	return ok
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *StorageRepository) Create(hostname, ip string, record dns.DNSResourceRecord) error {
	if !imsp.ExistsHostname(hostname) {
		imsp.Entries[hostname] = make(map[string]dns.DNSResourceRecord)
	}
	imsp.mutex.Lock()
	imsp.Entries[hostname][ip] = record
	imsp.LbInfo[hostname] = dns.NodeHealth{
		Healthy:     false,
		Connections: 0,
	}
	imsp.mutex.Unlock()
	if !imsp.Exists(hostname, ip) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}

// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *StorageRepository) Read(hostname string) (dns.DNSResourceRecord, error) {
	if !imsp.ExistsHostname(hostname) {
		return dns.DNSResourceRecord{}, errors.New("hostname does not exist")
	}
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()

	ip := findBestNode(hostname, imsp)

	return imsp.Entries[hostname][ip], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *StorageRepository) Update(hostname, ip string, value dns.DNSResourceRecord) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()

	// TODO: check hostname existence

	imsp.Entries[hostname][ip] = value
	imsp.LbInfo[ip] = dns.NodeHealth{
		Healthy:     false,
		Connections: 0,
	}
	return nil
}

// Adding the delete function, which basically just removes an specific entry (= the given service)
func (imsp *StorageRepository) Delete(hostname, ip string) error {
	if imsp.Exists(hostname, ip) {
		imsp.mutex.Lock()
		defer imsp.mutex.Unlock()
		delete(imsp.Entries[hostname], ip)
		if len(imsp.Entries[hostname]) == 0 {
			delete(imsp.Entries, hostname)
		}
		delete(imsp.LbInfo, ip)
	} else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong to delete the entry")
	}
	return nil
}

func findBestNode(hostname string, imsp *StorageRepository) string {
	lowestConnections := 99999
	var lowestIp string
	for ip := range imsp.Entries[hostname] {
		if imsp.LbInfo[ip].Connections < lowestConnections {
			lowestConnections = imsp.LbInfo[ip].Connections
			lowestIp = ip
		}
	}

	nodeHealth := imsp.LbInfo[lowestIp]
	nodeHealth.Connections += 1
	imsp.LbInfo[lowestIp] = nodeHealth

	go func() {
		time.Sleep(time.Duration(imsp.Entries[hostname][lowestIp].TimeToLive) * time.Second)
		if imsp.Exists(hostname, lowestIp) {
			nodeHealth := imsp.LbInfo[lowestIp]
			nodeHealth.Connections -= 1
			imsp.LbInfo[lowestIp] = nodeHealth
		}
	}()

	return lowestIp
}
