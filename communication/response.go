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
		"type": etype,
		"value": value,
	}
	response := struct {
		Service string `bson:"service"`
		Type string `bson:"type"`
		Key string `bson:"key"`
		Value map[string]string `bson:"value"`
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
		Domain string `bson:"domain"`
		Error string `bson:"error"`
		Value string `bson:"value"`
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
		"type": etype,
		"value": value,
	}
	response := struct {
		Service string `json:"service"`
		Type string `json:"type"`
		Key string `json:"key"`
		Value map[string]string `json:"value"`
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