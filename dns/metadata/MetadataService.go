package metadata

import (
	""
)

/*
This class is also seen as the "Model" it contains logical functionalities for the 
*/

// Predefined variables for the usage in this class
type repository = MetadataRepository
type requestType = string

type MetadataService interface {
	Create(repository MetadataRepository) MetadataServiceImpl
	HandleRequest(service, ip, key, requestType string) error
	POSTRequest(requestType string) bool
	GETRequest(requestType string) bool
}

type MetadataServiceImpl struct {
	repository MetadataRepository
}

func Create(repository MetadataRepository) MetadataServiceImpl {
	return MetadataServiceImpl{
		repository: repository,
	}
}

func (m *MetadataServiceImpl) ... {
	if m.repository.Exists(key) {
		...
	} else {
		err := m.repository.Create(...)
		if err != nil {
			panic(err)
		}
		...
	}
}