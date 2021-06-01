package communication

import (
	"bufio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var ip = net.ParseIP("127.0.0.1")
const port = 8000
const leaderaddr = "leader"

func TestMain(m *testing.M){
	zerolog.SetGlobalLevel(zerolog.Level(3))
	udpserver := UDPServer{
		Address: ip,
		Port:    port,
	}

	log.Info().Msg("Starting DNS")
	go udpserver.Start(func(addr net.Addr, buf []byte) []byte {
		return buf
	})

	time.Sleep(2 * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestUDPIsReachableShouldPass(t *testing.T) {
	conn, err := net.Dial("udp", ip.String()+":"+strconv.Itoa(port))
	//defer conn.Close()
	if err != nil {
		log.Error().Msgf("Error on establishing connection: %s\n", err)
	}
	msg := []byte("test")
	conn.Write(msg)
	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)

	if err != nil {
		t.Errorf("Test failed. Cannot read answer")
	}
}
func TestUDPIsReachableOverNetworkShouldPass(t *testing.T) {
	conn, err := net.Dial("udp", leaderaddr+":"+strconv.Itoa(port))
	//defer conn.Close()
	if err != nil {
		log.Error().Msgf("Error on establishing connection: %s\n", err)
	}
	msg := []byte("test")
	conn.Write(msg)
	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)

	if reflect.DeepEqual(answer, msg) {
		t.Errorf("Returned message (%s) does not match sent message (%s)",string(answer),string(msg))
	}
}

func TestUDPAnswersProperlyShouldPass(t *testing.T) {
	conn, err := net.Dial("udp", ip.String()+":"+strconv.Itoa(port))
	//defer conn.Close()
	if err != nil {
		log.Error().Msgf("Error on establishing connection: %s\n", err)
	}
	msg := []byte("test")
	conn.Write(msg)
	answer := make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(answer)

	if reflect.DeepEqual(answer, msg) {
		t.Errorf("Returned message (%s) does not match sent message (%s)",string(answer),string(msg))
	}
}