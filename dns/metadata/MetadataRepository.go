package metadata

/*
	This class is also seen as the "Repository" it contains the basic CRUD (= Create, Read, Update, Delete) functionalities
	for the metadata storage, which is nearly linked to the API.
*/

// Some necessary imports
import (
	"errors"
	"sync"
	"time"
)

// Creating the MetadataRepository interface which predefines the included functions of this class
type MetadataRepository interface {
	Exists(service, ip, key string) bool
	Create(service, ip, key, value string) error
	Read(service, ip, key string) (string, error)
	Update(service, ip, key, value string) error
	Delete(service, ip, key string) error
}

// Predefined variables for the following global one
type service = string
type ip = string
type key = string
type value = string

// Global metadata variable for this class
var metadata = make(map[service]map[ip]map[key]value)
var metadataMu = &sync.RWMutex{}


// Creating a structure for the metadata
type InMemoryStorageRepository struct {
	Location string `json:"Location"`
	Sensortype string `json:"Sensortype"`
	Registered time.Time `json:"Registered"`
	IsActive bool `json:"IsActive"`
}

// Adding the exists function, which basically just checks if an entry for this specific service exists
func (im *InMemoryStorageRepository) Exists(service, ip, key string) bool {
	if metadata[service][ip][key] != "" {
		return true
	}
	return false
}

// Adding the create function, which basically just adds a new entry to the map
func (im *InMemoryStorageRepository) Create(service, ip, key, value string) error {
	metadataMu.Lock()
	metadata[service][ip][key] = value
	metadataMu.Unlock()
	if !im.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return errors.New("success: a new entry was made")
}

func (im *InMemoryStorageRepository) Read(service, ip, key string) (string, error) {
	metadataMu.RLock()
	defer metadataMu.RUnlock()
	if !im.Exists(service, ip, key) {
		return metadata[service][ip][key], errors.New("no entry: as it looks like for this specific service no entry was made")
	}
	return metadata[service][ip][key], errors.New("success: there was an existing entry")
}

func (im *InMemoryStorageRepository) Update(service, ip, key, value string) error {
	metadataMu.Lock()
	metadata[service][ip][key] = value
	metadataMu.Unlock()
	if !im.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return errors.New("success: the entry was updated")
}

func (im *InMemoryStorageRepository) Delete(service, ip, key string) error {
	metadataMu.Lock()
	delete(metadata[service][ip], key)
	metadataMu.Unlock()
	if im.Exists(service, ip, key) {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return errors.New("success: the entry was deleted")
}