package metadata

/*
This class is also seen as the "Model" it contains logical functionalities for the 
*/

type MetadataService interface {
...
}

type MetadataServiceImpl struct {
	repository MetadataRepository
}
func Create(repository MetadataRepository) MetadataService {
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
