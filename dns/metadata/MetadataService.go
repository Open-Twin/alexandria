package metadata

import (
	""
	"errors"
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
	POSTRequest(requestType string) bool
	GETRequest(requestType string) bool
	ResponseRequest()
}

type MetadataServiceImpl struct {
	repository MetadataRepository
}

func Init(repository MetadataRepository) MetadataServiceImpl {
	return MetadataServiceImpl{
		repository: repository,
	}
}

func POSTRequest(requestType string) bool {
	...
}

func GETRequest(requestType string) bool {
	...
}

func ResponseRequest() {
	...
}

func (m *MetadataServiceImpl) HandleRequest (service, ip, key, value, requestType string) error {
	if !m.repository.Exists(service, ip, key) {
		if requestType == "get" {
			return errors.New("not existing: the given key doesnt exist, thus no data can be returned")
		}
		err := m.repository.Create(service, ip, key, value)
		if err != nil {
			panic(err)
		}

	}
}