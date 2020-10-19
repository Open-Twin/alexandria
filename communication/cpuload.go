package communication

import (
	"log"
	"net/http"
	"time"
)

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) cpuLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tpe", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("LMAO TEST"))
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		h.logger.Println("request processed in %s\n", time.Now().Sub(startTime))
	}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/cpuload", h.Logger(h.cpuLoad))
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}
