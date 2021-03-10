package storage

import (
	"errors"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/loadbalancing"
	"sync"
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
	Entries map[hostname]map[ip]record
	// Creating a mutex onto the Metadata variable in order to handle threads
	mutex  sync.RWMutex
	LbInfo map[ip]loadbalancing.NodeHealth
}

func NewInMemoryDNSStorageRepository() *StorageRepository {
	entries := make(map[hostname]map[ip]record)
	info := make(map[ip]loadbalancing.NodeHealth)
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
	imsp.LbInfo[hostname] = loadbalancing.NodeHealth{
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
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()

	ip := loadbalancing.FindBestNode(hostname, imsp)

	return imsp.Entries[hostname][ip], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *StorageRepository) Update(hostname, ip string, value dns.DNSResourceRecord) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	imsp.Entries[hostname][ip] = value
	imsp.LbInfo[ip] = loadbalancing.NodeHealth{
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
		delete(imsp.LbInfo, ip)
	} else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong, to delete the entry")
	}
	return nil
}
