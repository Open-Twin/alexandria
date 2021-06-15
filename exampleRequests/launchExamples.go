package main

import (
	"bufio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(1)

	if os.Args[1] == "dns" {
		sendDnsEntry()
	} else if os.Args[1] == "meta" {
		sendMetadataPacket()
	}
}

func sendBsonMessage(address string, msg bson.M) {
	conn, err := net.Dial("udp", address)
	defer conn.Close()
	if err != nil {
		log.Error().Msgf("Error on establishing connection: %s", err)
	}
	sendMsg, err := bson.Marshal(msg)

	conn.Write(sendMsg)
	log.Info().Msgf("Message sent: %s", sendMsg)

	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)
	if err != nil {
		log.Error().Msgf("Error on receiving answer: %v", err)
	} else {
		log.Info().Msgf("Answer:\n %s", answer)
	}
}

func sendDnsEntry() {
	address := "192.168.198.128:10000"
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
		"hostname":    os.Args[2],
		"ip":          os.Args[3],
		"requestType": os.Args[4],
	}

	sendBsonMessage(address, msg)
}

func sendMetadataPacket() {
	address := "192.168.198.128:20000"
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
		"service": "toilet",
		"type":    os.Args[2],
		"key":     "temp",
		"value":   os.Args[3],
	}

	sendBsonMessage(address, msg)
}
