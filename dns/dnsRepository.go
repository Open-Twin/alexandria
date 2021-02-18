package dns

import (
	"errors"
	"sync"
)

/*
	This class is also seen as the "Repository" it contains the basic CRUD (= Create, Read, Update, Delete) functionalities
	for the Metadata storage, which is nearly linked to the API.
*/

// Creating the MetadataRepository interface which predefines the needed methods for the storage
type MetadataRepository interface {
	Exists(service, ip, key string) bool
	Create(service, ip, key, value string) error
	Read(service, ip, key string) (string, error)
	Update(service, ip, key, value string) error
	Delete(service, ip, key string) error
}

// Predefined variables for the usage in this class
type hostname = string
type rdata = string
type zone = string
type record = string

// Creating a structure for the Metadata containing the necessary variables
type StorageRepository struct {
	// Global Metadata variable for this class
	Entries map[hostname]map[record]map[zone]rdata
	// Creating a mutex onto the Metadata variable in order to handle threads
	mutex sync.RWMutex
}
func NewInMemoryStorageRepository() *StorageRepository{
	entries := make(map[hostname]map[record]map[zone]rdata)
	return &StorageRepository{
		Entries: entries,
	}
}
// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *StorageRepository) Exists(service, ip, key string) bool {
	if imsp.Entries[service][ip][key] != "" {
		return true
	}

	return false
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *StorageRepository) Create(service, ip, key, value string) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()



	if imsp.Exists(service, ip, key) {
		imsp.Entries[service][ip][key] = value
	} else {
		imsp.Entries[service]=make(map[string]map[string]string)
		imsp.Entries[service][ip] = make(map[string]string)
		imsp.Entries[service][ip][key] = value
		if !imsp.Exists(service, ip, key) {
			return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
		}
	}
	return nil
}


// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *StorageRepository) Read(service, ip, key string) (string, error) {
	imsp.mutex.RLock()
	defer imsp.mutex.RUnlock()
	if !imsp.Exists(service, ip, key) {
		return imsp.Entries[service][ip][key], errors.New("no entry: as it looks like for this specific service no entry was made")
	}
	return imsp.Entries[service][ip][key], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *StorageRepository) Update(service, ip, key, value string) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	if imsp.Exists(service, ip, key) {
		imsp.Entries[service][ip][key] = value
	}else {
		imsp.Entries[service][ip][key] = value
		if !imsp.Exists(service, ip, key) {
			return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
		}
	}
	return nil
}

// Adding the delete function, which basically just removes an specific entry (= the given service)
func (imsp *StorageRepository) Delete(service, ip, key string) error {
	imsp.mutex.Lock()
	defer imsp.mutex.Unlock()
	if imsp.Exists(service, ip, key) {
		delete(imsp.Entries[service][ip], key)
	}else {
		delete(imsp.Entries[service][ip], key)
		//! hinzugefügt
		if !imsp.Exists(service, ip, key) {
			return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
		}
	}
	return nil
}