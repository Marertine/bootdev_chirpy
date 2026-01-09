package main

import (
	"net/http"
)

func handlerReset(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}
