package main

import (
	"github.com/Open-Twin/alexandria/communication"
	"log"
	"net/http"
	"os"
)

var (
	certFile   = ""
	keyFile    = ""
	serverAddr = ""
)

func main() {
	logger := log.New(os.Stdout, "healthcheck", log.LstdFlags|log.Lshortfile)

	cpu := communication.NewHandlers(logger)

	mux := http.NewServeMux()
	cpu.SetupRoutes(mux)

	srv := communication.New(mux)

	logger.Println("server gestaretet")
	err := srv.ListenAndServe()
	//err := srv.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
