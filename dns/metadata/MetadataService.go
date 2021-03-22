package metadata

import (
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
	POSTRequest(repository MetadataRepository, service, ip, key, value string) error
	GETRequest(repository MetadataRepository, service, ip, key string) error
	UPDATERequest(repository MetadataRepository, service, ip, key, value string) error
	DELETERequest(repository MetadataRepository, service, ip, key string) error
	ResponseRequest(repository MetadataRepository) error
}

type MetadataServiceImpl struct {
	repository MetadataRepository
}

func Init(repository MetadataRepository) MetadataServiceImpl {
	return MetadataServiceImpl{
		repository: repository,
	}
}

func POSTRequest(repository MetadataRepository, service, ip, key, value string) error {
	err := repository.Create(service, ip, key, value)
	if err != nil {
		return err
	}
	return nil
}

func GETRequest(repository MetadataRepository, service, ip, key string) error {
	_, err := repository.Read(service, ip, key)
	if err != nil {
		return err
	}
	return nil
}

func UPDATERequest(repository MetadataRepository, service, ip, key, value string) error {
	err := repository.Update(service, ip, key, value)
	if err != nil {
		return err
	}
	return nil
}

func DELETERequest(repository MetadataRepository, service, ip, key string) error {
	err := repository.Delete(service, ip, key)
	if err != nil {
		return err
	}
	return nil
}

func ResponseRequest(repository MetadataRepository, service, ip, key, value, requestType string, err error) error {
	if err != nil {
		
	}
}

func (msi *MetadataServiceImpl) HandleRequest (service, ip, key, value, requestType string) error {
	if !msi.repository.Exists(service, ip, key) {
		switch requestType {
			case "get":
				return errors.New("not existing: the given key doesnt exist, thus no data can be returned")
			case "update":
				return errors.New("not existing: the given key doesnt exist, thus no data can be updated")
			case "delete":
				return errors.New("not existing: the given key doesnt exist, thus no data can be deleted")
			case "store":
				err := POSTRequest(msi.repository, service, ip, key, value)
				if err != nil {
					return err
				}
				return nil
			default:
				return errors.New("unsupported operation: the given operation you'd like to execute isn't supported")
		}
	}
	switch requestType {
		case "get":
			err := GETRequest(msi.repository, service, ip, key)
			if err != nil {
				return err
			}
			return nil
		case "update":
			err := UPDATERequest(msi.repository, service, ip, key, value)
			if err != nil {
				return err
			}
			return nil
		case "delete":
			err := DELETERequest(msi.repository, service, ip, key)
			if err != nil {
				return err
			}
			return nil
		case "store":
			err := POSTRequest(msi.repository, service, ip, key, value)
			if err != nil {
				return err
			}
			return nil
		default:
			return errors.New("unsupported operation: the given operation you'd like to execute isn't supported")
	}
}