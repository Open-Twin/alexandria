package metadata

/*
	This class is there to represent the basic storage backend of for the metadata which is received from the services
	of the service mesh system.
*/

import (
	"encoding/json"
	"fmt"
	"time"
)

// Creating a struct for the services, in which the metadata struct will be located so it is easier afterwards to remove a certain service with its data
type Service []struct {
	ID string
	Data Metadata
}

// Creating a structure for the metadata json, which will be later on encoded to json via marshal
type Metadata struct {
	Location string
	Sensortype string
	RequestReceived time.Time
	IsActive bool
}

// Creating a function to convert the struct into a json
func StructToJson(data Metadata) []byte {
	jsondata, _ := json.Marshal(data)
	return jsondata
}

// Creating a function to return a specified service from the service array struct
func getService() {
	// To be continued
}

// Creating a function to save a specified service into the service array struct
func setSerService() {
	// To be continued
}

// Creating a function to delete a specified service from the service array struct
func deleteService() {
	// To be continued
}

// Creating a "main" function to test the functionality of the storage class so far
func Main() {
	//define a dummy struct
	data := Metadata{
		Location: "Vienna",
		Sensortype: "Water",
		RequestReceived: time.Now(),
		IsActive: true,
	}

	// encode the data dummy into a json
	jsondata := StructToJson(data)

	// printing the json as a string into the console
	fmt.Println(string(jsondata))
}