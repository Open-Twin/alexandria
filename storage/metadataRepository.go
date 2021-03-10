package storage

/*
	This class is also seen as the "Repository" it contains the basic CRUD (= Create, Read, Update, Delete) functionalities
	for the Metadata storage, which is nearly linked to the API.
*/

// Some necessary imports, that are needed for the storage backend
import (
	"errors"
	"sync"
)

// Creating the MetadataRepository interface which predefines the needed methods for the storage
type MetadataRepository interface {
	Exists(service, ip, key string) bool
	Create(service, ip, key, value string) error
	Read(service, ip, key string) (string, error)
	Update(service, ip, key, value string) error
	Delete(service, ip, key string) error
}

// Creating a structure for the Metadata containing the necessary variables
type InMemoryStorageRepository struct {
	// Global Metadata variable for this class
	Metadata map[service]map[ip]map[key]value
	// Creating a mutex onto the Metadata variable in order to handle threads
	metadataMu sync.RWMutex
}

func NewInMemoryStorageRepository() *InMemoryStorageRepository {
	 metadata := make(map[service]map[ip]map[key]value)
	 return &InMemoryStorageRepository{
	 	Metadata: metadata,
	 }
}
// Adding the exists function, which basically just checks if an entry for this specific service exists
func (imsp *InMemoryStorageRepository) Exists(service, ip, key string) bool {
	imsp.metadataMu.RLock()
	defer imsp.metadataMu.RUnlock()
	_, ok := imsp.Metadata[service][ip][key]
	return ok
}

// Adding the create function, which basically just adds a new entry to the map
func (imsp *InMemoryStorageRepository) Create(service, ip, key, value string) error {
	if imsp.Exists(service, ip, key) {
		imsp.metadataMu.Lock()
		imsp.Metadata[service][ip][key] = value
		imsp.metadataMu.Unlock()
	} else {
		imsp.metadataMu.Lock()
		imsp.Metadata[service]=make(map[string]map[string]string)
		imsp.Metadata[service][ip] = make(map[string]string)
		imsp.Metadata[service][ip][key] = value
		imsp.metadataMu.Unlock()
		if !imsp.Exists(service, ip, key) {
			return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
		}
	}
	return nil
}


// Adding the read function, which basically just returns the specific value of the given service as a string
func (imsp *InMemoryStorageRepository) Read(service, ip, key string) (string, error) {

	defer imsp.metadataMu.RUnlock()

	// TODO: Loadbalancer

	if !imsp.Exists(service, ip, key) {
		imsp.metadataMu.RLock()
		return imsp.Metadata[service][ip][key], errors.New("no entry: as it looks like for this specific service no entry was made")
	}
	imsp.metadataMu.RLock()
	return imsp.Metadata[service][ip][key], nil
}

// Adding the update function, which basically just replaces a specific value of the given service with the new given value
func (imsp *InMemoryStorageRepository) Update(service, ip, key, value string) error {
	if imsp.Exists(service, ip, key) {
		imsp.metadataMu.Lock()
		imsp.Metadata[service][ip][key] = value
		defer imsp.metadataMu.Unlock()
	}else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}

// Adding the delete function, which basically just removes a specific entry (= the given service)
func (imsp *InMemoryStorageRepository) Delete(service, ip, key string) error {
	if imsp.Exists(service, ip, key) {
		imsp.metadataMu.Lock()
		delete(imsp.Metadata[service][ip], key)
		imsp.metadataMu.Unlock()
	}else {
		return errors.New("wrong argument: probably one of the given arguments is either non existing or wrong")
	}
	return nil
}