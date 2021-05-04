package dns

/*import (
	"bufio"
	"errors"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

func RecursiveLookup(originalMessage []byte) ([]byte, error){
	//send request to outside dns from conf file
	conf := cfg.ReadConf()
	resultChannel := make(chan []byte)
	for _, rec := range conf.NSRecords{
		go sendRecursiveRequest(rec, originalMessage, resultChannel)
	}
	//TODO: timeout and check if result is good?
	result := <- resultChannel
	return result, nil

	//return result, errors.New("recursive query had no success")
}

func sendRecursiveRequest(record string, buf []byte, channel chan []byte){
	// Setup a UDP connection
	conn, err := net.Dial("udp", record)
	if err != nil {
		log.Fatal().Msg("failed to connect:", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal().Msg("failed to set deadline: ", err)
	}
	conn.Write(buf)

	encodedAnswer := make([]byte, len(buf))
	if _, err := bufio.NewReader(conn).Read(encodedAnswer); err != nil {
		//return nil,err
		return
	}
	channel <- encodedAnswer
	return
	//return encodedAnswer, nil
}*/