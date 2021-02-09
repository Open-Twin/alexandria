package loadbalancing

import (
	"fmt"
	"net/http"
	"time"
)

func StartLoadReporting() {
	http.HandleFunc("/health", sendLoad)
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Started reporting current server load")
}

func sendLoad(w http.ResponseWriter, r *http.Request) {
	data := collectData()
	fmt.Printf("Sending status as requested: %s\n", string(data))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func collectData() []byte {
	timestamp := time.Now()
	return []byte(timestamp.String())
}
