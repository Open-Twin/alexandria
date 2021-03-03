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

// Predefined variables for the usage in this class

type hostname = string
type values = []dns.DNSResourceRecord

// Creating a structure for the Metadata containing the necessary variables
type StorageRepository struct {
	Entries map[hostname]values
	// Creating a mutex onto the Metadata variable in order to handle threads
	mutex sync.RWMutex
}

func NewInMemoryDNSStorageRepository() *StorageRepository {
	entries := make(map[hostname]values)
	return &StorageRepository{
		Entries: entries,
	}
}
// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *StorageRepository) Exists(hostname string) bool {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	_, ok := imsp.Entries[hostname]
	return ok
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *StorageRepository) Create(hostname string, record dns.DNSResourceRecord) error {
	log.Print("CREATE: "+hostname)
	//imsp.Entries[hostname]=make(map[string]dns.DNSResourceRecord)
	imsp.mutex.Lock()
	imsp.Entries[hostname] = append(imsp.Entries[hostname], record)
	imsp.mutex.Unlock()
	if !imsp.Exists(hostname) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}

// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *StorageRepository) Read(hostname string) (dns.DNSResourceRecord, error) {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	if !imsp.Exists(hostname) {
		//TODO: works?
		return imsp.Entries[hostname][0], errors.New("no entry: as it looks like for this specific service no entry was made")
	}
	//TODO: query with loadbalacing algorithms (least connected?)
	log.Print("READ: "+hostname)
	log.Println(imsp.Entries[hostname])
	log.Print("----")
	return imsp.Entries[hostname][0], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *StorageRepository) Update(hostname string, value dns.DNSResourceRecord) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	/*if imsp.Exists(hostname) {
		imsp.Entries[hostname] = value
	}else {
		return errors.New("not existing: the entry you'd like to update didnt exist, instead it was created")
	}*/
	//TODO: langsam?
	index := getIndexFromResourceRecord(hostname,imsp.Entries)
	if index == -1 {
		return errors.New("not existing: the entry you'd like to update didnt exist")
	}
	imsp.Entries[hostname][index] = value
	return nil
}

// Adding the delete function, which basically just removes an specific entry (= the given service)
func (imsp *StorageRepository) Delete(hostname string) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	/*if imsp.Exists(hostname) {
		delete(imsp.Entries, hostname)
	}else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong, to delete the entry")
	}*/
	//TODO: langsam?
	index := getIndexFromResourceRecord(hostname,imsp.Entries)
	if index == -1 {
		return errors.New("not existing: the entry you'd like to update didnt exist")
	}
	//delete from slice
	//TODO: langsam?
	imsp.Entries[hostname][index] = imsp.Entries[hostname][len(imsp.Entries[hostname])-1]
	newSlice := imsp.Entries[hostname][:len(imsp.Entries[hostname])-1]
	imsp.Entries[hostname] = newSlice
	return nil
}

func getIndexFromResourceRecord(hostname string, m map[hostname]values) int{
	for i, obj := range m[hostname] {
		domainName := ""
		for i2, part := range obj.Labels {
			domainName += part
			if i2 < len(obj.Labels)-1 {
				domainName += "."
			}
		}
		if domainName == hostname{
			return i
		}
	}
	return -1
}