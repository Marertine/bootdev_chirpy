package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Boot.Dev/Twitter Clone")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

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
