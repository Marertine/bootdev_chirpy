package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Boot.Dev/Twitter Clone")

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthzHandler)

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

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
