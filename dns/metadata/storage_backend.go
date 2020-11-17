package metadata

/*
	This class is there to represent the basic storage backend of for the metadata which is received from the services
	of the service mesh system.
*/

// Some necessary imports
import (
	"encoding/json"
	"fmt"
	"time"
)

// Global variable which is a map and contains all services
var services = make(map[string]Metadata)

// Creating a structure for the metadata json, which will be later on encoded to json via marshal
type Metadata struct {
	Location string `json:"Location"`
	Sensortype string `json:"Sensortype"`
	Registered time.Time `json:"Registered"`
	IsActive bool `json:"IsActive"`
}

// Creating a function to convert the struct into a json
func StructToJson(data Metadata) []byte {
	jsondata, _ := json.Marshal(data)
	return jsondata
}

// Creating a function to convert a json into the struct
func JsonToStruct(jsonData string) Metadata {
	data := Metadata{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Following error occurred:", err)
	}
	return data
}

// Creating a function to return the metadata of a specified service from the global service map
func getServiceMeta(ip string) Metadata{
	return services[ip]
}

// Creating a function to return the metadata of a specified service from the global service map as a Json
func getServiceJson(ip string) string{
	json := StructToJson(services[ip])
	return string(json)
}

// Creating a function to change the data of a specified service from the global service map
func changeService(ip string, data Metadata) {
	services[ip] = data
}

// Creating a function to change the data of a specified service from the global service map
func addService(ip string, data Metadata) {
	services[ip] = data
}

// Creating a function to delete a specified service from the service array struct
func deleteService(ip string) {
	delete(services,ip)
}

// Creating a "main" function to test the functionality of the storage class so far
func Main() {
	//define a dummy struct
	wataData := Metadata{
		Location: "Vienna",
		Sensortype: "Water",
		Registered: time.Now(),
		IsActive: false,
	}

	jsonWata := StructToJson(wataData)
	fmt.Println(string(jsonWata))
}