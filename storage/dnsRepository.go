package storage

import (
	"errors"
	"github.com/Open-Twin/alexandria/dns"
	"log"
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
	mutex sync.RWMutex
}

func NewInMemoryDNSStorageRepository() *StorageRepository {
	entries := make(map[hostname]map[ip]record)
	return &StorageRepository{
		Entries: entries,
	}
}
// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *StorageRepository) Exists(hostname, ip string) bool {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	_, ok := imsp.Entries[hostname][ip]
	return ok
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *StorageRepository) Create(hostname, ip string, record dns.DNSResourceRecord) error {
	log.Print("CREATE: "+hostname)
	imsp.mutex.Lock()
	imsp.Entries[hostname]=make(map[string]dns.DNSResourceRecord)
	imsp.Entries[hostname][ip] = record
	imsp.mutex.Unlock()
	if !imsp.Exists(hostname,ip) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}

// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *StorageRepository) Read(hostname string) (dns.DNSResourceRecord, error) {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	//TODO: query with loadbalacing algorithms (least connected?)
	/*if !imsp.Exists(hostname) {
		return imsp.Entries[hostname][ip], errors.New("no entry: as it looks like for this specific service no entry was made")
	}*/
	log.Print("READ: "+hostname)
	//only temporary for testing
	for k := range imsp.Entries[hostname]{
		log.Println(imsp.Entries[hostname][k])
		log.Print("----")
		return imsp.Entries[hostname][k], nil
	}
	return imsp.Entries[hostname][""], nil
	//return imsp.Entries[hostname][ip], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *StorageRepository) Update(hostname, ip string, value dns.DNSResourceRecord) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	imsp.Entries[hostname][ip] = value
	return nil
}

// Adding the delete function, which basically just removes an specific entry (= the given service)
func (imsp *StorageRepository) Delete(hostname, ip string) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	if imsp.Exists(hostname,ip) {
		delete(imsp.Entries[hostname], ip)
	}else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong, to delete the entry")
	}
	return nil
}