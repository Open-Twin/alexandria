package main

import (
	"gopkg.in/mgo.v2/bson"
	"os"
)

func main() {
	address := "127.0.0.1:10000"
	/*msg := bson.M{
		"Labels":             []string{"at", "ac", "dejan"},
		"Type":               uint16(60),
		"Class":              uint16(50),
		"TimeToLive":         uint32(32),
		"ResourceDataLength": uint16(30),
		"ResourceData":       []byte{0x16, 0x32, 0x64},
		"RequestType":		  "store",
	}*/
	msg := bson.M{
		"hostname":    os.Args[1],
		"ip":          os.Args[2],
		"requestType": "store",
	}

	sendBsonMessage(address, msg)
}
