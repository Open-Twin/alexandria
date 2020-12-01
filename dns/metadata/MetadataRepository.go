package metadata

/*
	This class is also seen as the "Repository" it contains the basic CRUD (= Create, Read, Update, Delete) functionalities
	for the metadata storage, which is nearly linked to the API.
*/

// Some necessary imports
import (
	"errors"
	"sync"
)

// Creating the MetadataRepository interface which predefines the included functions of this class
type MetadataRepository interface {
	Exists(service, ip, key string) bool
	Create(service, ip, key, value string) error
	Read(service, ip, key string) (string, error)
	Update(service, ip, key, value string) error
	Delete(service, ip, key string) error
}

// Predefined variables for the usage in this class
type service = string
type ip = string
type key = string
type value = string

// Creating a structure for the metadata containing the necessary variables
type InMemoryStorageRepository struct {
	// Global metadata variable for this class
	metadata map[service]map[ip]map[key]value
	metadataMu sync.RWMutex
}

// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *InMemoryStorageRepository) Exists(service, ip, key string) bool {
	if imsp.metadata[service][ip][key] != "" {
		return true
	}
	return false
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *InMemoryStorageRepository) Create(service, ip, key, value string) error {
	imsp.metadataMu.Lock()
	defer imsp.metadataMu.Unlock()
	imsp.metadata[service][ip][key] = value
	if !imsp.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}


// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *InMemoryStorageRepository) Read(service, ip, key string) (string, error) {
	imsp.metadataMu.RLock()
	defer imsp.metadataMu.RUnlock()
	if !imsp.Exists(service, ip, key) {
		return imsp.metadata[service][ip][key], errors.New("no entry: as it looks like for this specific service no entry was made")
	}
	return imsp.metadata[service][ip][key], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *InMemoryStorageRepository) Update(service, ip, key, value string) error {
	imsp.metadataMu.Lock()
	defer imsp.metadataMu.Unlock()
	imsp.metadata[service][ip][key] = value
	if !imsp.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}

// Adding the delete function, which basically just removes an specific entry (= the given service)
func (imsp *InMemoryStorageRepository) Delete(service, ip, key string) error {
	imsp.metadataMu.Lock()
	defer imsp.metadataMu.Unlock()
	delete(imsp.metadata[service][ip], key)
	if imsp.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}