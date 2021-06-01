package loadbalancing

import (
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"time"
)

func StartLoadReporting(lbUrl string) {
	lbRegister(lbUrl)

	http.HandleFunc("/health", sendLoad)
	//TODO: PORTS

	go http.ListenAndServe(":8080", nil)
	log.Info().Msg("Started reporting alexandria server load")
}

func lbRegister(lbUrl string) {
	//TODO: PORTS
	resp, err := http.Get("http://" + lbUrl + ":8080/signup")
	if err != nil {
		log.Error().Msgf("Registration at loadbalancer failed: %v", err)
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("Error reading loadbalancer answer: %v", err)
		}

		if string(body) != "succesfully added" {
			log.Error().Msgf("Adding node didn't work: %v", string(body))
		}

		log.Info().Msgf("Registered at loadbalancer")
	}
}

func sendLoad(w http.ResponseWriter, r *http.Request) {
	data := collectData()
	log.Debug().Msgf("Sending status as requested: %s\n", string(data))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func collectData() []byte {
	timestamp := time.Now()
	return []byte(timestamp.String())
}
