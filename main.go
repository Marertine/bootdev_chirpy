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

	// Admin handlers
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		handlerReset(w, r, &apiCfg)
	})

	// Application handlers
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	// API handlers
	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		handlerMetrics(w, r, &apiCfg)
	})

	mux.HandleFunc("POST /api/validate_chirp", handlerValidation)

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
