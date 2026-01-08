package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	fmt.Println("Boot.Dev/Twitter Clone")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	filepathRoot := "."
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r, &apiCfg)
	})

	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		resetHandler(w, r, &apiCfg)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the HTTP server
	log.Println("Starting server on :8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func resetHandler(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}
