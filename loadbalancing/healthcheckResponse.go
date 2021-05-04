package loadbalancing

import (
	"io/ioutil"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"time"
)

func StartLoadReporting(lbUrl string) {
	lbRegister(lbUrl)

	http.HandleFunc("/health", sendLoad)
	go http.ListenAndServe(":8080", nil)
	log.Info().Msg("Started reporting alexandria server load")
}

func lbRegister(lbUrl string) {
	data := url.Values{
		"ip": {"127.0.0.1"},
	}

	resp, err := http.PostForm("http://"+lbUrl+"/signup", data)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
	}

	if string(body) != "succesfully added" {
		log.Error().Msgf("Adding node didn't work: %v", string(body))
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
