package loadbalancing

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func StartLoadReporting() {
	http.HandleFunc("/health", sendLoad)
	go http.ListenAndServe(":8080", nil)
	log.Info().Msg("Started reporting alexandria server load")
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
