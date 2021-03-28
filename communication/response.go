package communication

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

/*
 * create response for metadata requests
 */
func CreateMetadataResponse(service, key, etype, value string) []byte{

	valueMap := map[string]string{
		"Type": etype,
		"Value": value,
	}
	response := struct {
		Service string `bson:"Service"`
		Type string `bson:"Type"`
		Key string `bson:"Key"`
		Value map[string]string `bson:"Value"`
	}{
		Service: service,
		Type: "response",
		Key: key,
		Value: valueMap,
	}

	responseBytes, err := bson.Marshal(response)
	if err != nil {
		log.Print("sendresponse failed")
	}

	return responseBytes
}

/*
 * create response for DNS requests
 */
func CreateResponse(domain, etype, value string) []byte{

	response := struct {
		Domain string
		Error string
		Value string
	}{
		Domain: domain,
		Error: etype,
		Value: value,
	}

	responseBytes, err := bson.Marshal(response)
	if err != nil {
		log.Print("sendresponse failed")
	}

	return responseBytes
}

/*
 * send http response to responsewriter
 */
func sendHttpResponse(service, key, etype, value string, w http.ResponseWriter){

	valueMap := map[string]string{
		"Type": etype,
		"Value": value,
	}
	response := struct {
		Service string
		Type string
		Key string
		Value map[string]string
	}{
		Service: service,
		Type: "response",
		Key: key,
		Value: valueMap,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		//server.Logger.Println("")
		log.Print("sendresponse failed")
	}

	w.Write(responseBytes)
}