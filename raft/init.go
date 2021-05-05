package raft

import (
	"errors"
	"fmt"
	"github.com/Open-Twin/alexandria/cfg"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Start(conf *cfg.Config) (*Node,error){

	raftNode, err := NewNode(conf)
	if err != nil {
		return nil, errors.New("Error configuring node: "+err.Error())
	}
	//attempts to join a node if join Address is given
	//if joinaddress is given, join with that address
	joinAddr := conf.JoinAddr
	raftAddr := conf.RaftAddr
	if conf.JoinAddr != nil {
		log.Info().Msg("attempting join")
		go join(joinAddr.String(), raftAddr.String(), 0)
	}else if conf.Autojoin {
		//else try to autojoin
		log.Info().Msg("attempting auto-joining")
		err := tryAutoJoin(raftAddr.String(), strconv.Itoa(conf.UdpPort), strconv.Itoa(conf.HttpAddr.Port))
		if err != nil{
			log.Error().Msg("autojoin failed: "+err.Error())
		}
	}
	return raftNode, nil
}

/*
 * tries to join node into cluster maxTries times.
 * if maxTries is 0, tries indefinitely.
 */
func join(joinaddr, raftaddr string, maxTries int) error {
	retryJoin := func() error {
		joinUrl := url.URL{
			Scheme: "http",
			Host:   joinaddr,
			Path:   "join",
		}
		req, err := http.NewRequest(http.MethodPost, joinUrl.String(), nil)
		req.Close = true
		if err != nil {
			return err
		}
		req.Header.Add("Peer-Address", raftaddr)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("non 200 status code: %d", resp.StatusCode)
		}
		log.Info().Msg("join successful: "+joinaddr)
		return nil
	}
	failedJoins := -1
	for {
		if err := retryJoin(); err != nil {
			log.Error().Msg("error joining cluster: "+joinaddr)
			failedJoins++
			if maxTries > 0 && failedJoins >= maxTries {
				return errors.New("exceeded maximum join tries")
			}
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}
}
/*
 * broadcasts a udp message to everyone in the network and listens
 * for a message from a fellow server. Tries to join said server after receiving message.
 */
func tryAutoJoin(raftaddr, broadcastPort, httpPort string) error {
	//broadcast udp to find available servers
	broadcastAddress := "255.255.255.255"
	pc, err := net.ListenPacket("udp4", ":"+broadcastPort)
	if err != nil {
		return err
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", broadcastAddress+":"+broadcastPort)
	if err != nil {
		return err
	}
	//broadcast udp message
	_, err = pc.WriteTo([]byte("autojoin-request"), addr)
	if err != nil {
		return err
	}
	//read responses
	for {
		buf := make([]byte, 1024)
		n, respaddr, err := pc.ReadFrom(buf)
		if err != nil {
			return err
		}
		log.Info().Msgf("autojoin received response from \"%s: %s\n", respaddr, buf[:n])
		//try joining address but with http port
		err = join(strings.Split(respaddr.String(),":")[0]+":"+httpPort, raftaddr, 5)
		if err != nil {
			log.Info().Msgf("failed autojoin request to "+respaddr.String()+":"+err.Error())
			log.Info().Msgf("waiting for other replies")
		}else{
			return nil
		}
	}

}


