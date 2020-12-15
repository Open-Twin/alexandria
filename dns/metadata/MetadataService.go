package metadata

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

/*
This class is also seen as the "Model" it contains logical functionalities for the 
*/

// Predefined variables for the usage in this class
type repository = MetadataRepository
type requestType = string

type MetadataService interface {
	Init(repository MetadataRepository) MetadataServiceImpl
	HandleRequest(service, ip, key, requestType string) error
	POSTRequest(repository MetadataRepository, service, ip, key, value string) bool
	GETRequest(repository MetadataRepository, service, ip, key string) bool
	UPDATERequest(repository MetadataRepository, service, ip, key, value string) bool
	DELETERequest(repository MetadataRepository, service, ip, key string) bool
	ResponseRequest(repository MetadataRepository) bool
}

type MetadataServiceImpl struct {
	repository MetadataRepository
}

func Init(repository MetadataRepository) MetadataServiceImpl {
	return MetadataServiceImpl{
		repository: repository,
	}
}

func POSTRequest(repository MetadataRepository, service, ip, key, value string) bool {

}

func GETRequest(repository MetadataRepository, service, ip, key string) bool {

}

func UPDATERequest(repository MetadataRepository, service, ip, key, value string) bool {

}

func DELETERequest(repository MetadataRepository, service, ip, key string) bool {

}

func ResponseRequest(repository MetadataRepository) bool {

}

func (msi *MetadataServiceImpl) HandleRequest (service, ip, key, value, requestType string) error {
	if !msi.repository.Exists(service, ip, key) {
		if requestType == "get" {
			return errors.New("not existing: the given key doesnt exist, thus no data can be returned")
		} else if requestType == "restore" {
			err := msi.repository.Create(service, ip, key, value)
			if err != nil {
				panic(err)
			}
		} else {
			return errors.New("unsupported operation: the given operation you'd like to execute isn't supported")
		}
	}
}