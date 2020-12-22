package loadbalancing

import (
	"fmt"
	"math/rand"
	"net/http"
)

func startLoadReporting() {
	http.HandleFunc("/load", sendLoad)
	fmt.Println("Started reporting current server load")
}

func sendLoad(w http.ResponseWriter, r *http.Request) {
	w.Write(collectData())
}

func collectData() []byte {
	randLoad := string(rand.Intn(99))
	return []byte(randLoad)
}
