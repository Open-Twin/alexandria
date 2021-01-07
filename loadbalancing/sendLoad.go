package loadbalancing

import (
	"fmt"
	"math/rand"
	"net/http"
)

func StartLoadReporting() {
	http.HandleFunc("/load", sendLoad)
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Started reporting current server load")
}

func sendLoad(w http.ResponseWriter, r *http.Request) {
	w.Write(collectData())
}

func collectData() []byte {
	randLoad := fmt.Sprint(rand.Intn(99))
	return []byte(randLoad)
}
