package metadata_test

import (
	"github.com/Open-Twin/alexandria/dns/metadata"
	"testing"
	"time"
)

// Global variable which is a map and contains all services
var services = make(map[string]metadata.Metadata)

func TestMethodStructToJson(t *testing.T) {
	waitAMin := time.Date(2021,time.January,1,1,1,1,1,time.UTC)

	validData := metadata.Metadata{
		Location: "Vienna",
		Sensortype: "Water",
		Registered: waitAMin,
		IsActive: true,
	}

	resultData := string(metadata.StructToJson(validData))
	expectedData := "{\"Location\":\"Vienna\",\"Sensortype\":\"Water\",\"Registered\":\"2021-01-01T01:01:01.000000001Z\",\"IsActive\":true}"
	if resultData != expectedData {
		t.Errorf("metadata.StructToJson failed, cause %q was expected as the result but instead got %q", expectedData, resultData)
	}
}

func TestValidJsonToStruct(t *testing.T) {

}

func TestInvalidJsonToStruct(t *testing.T) {

}

func TestGetServiceMeta(t *testing.T) {
	
}

func TestGetServiceJson(t *testing.T)  {

}

func TestChangeService(t *testing.T)  {

}

func TestAddService(t *testing.T) {

}

func TestDeleteService(t *testing.T) {

}