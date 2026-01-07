package main

import (
	"fmt"
	"log"
	"net/http"
)

type ServeMux struct {
	mux *http.ServeMux
}

func main() {
	fmt.Println("Boot.Dev/Twitter Clone")

	mux := NewServeMux()

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux.mux); err != nil {
		log.Fatal(err)
	}
}

// NewServeMux allocates and returns a new [ServeMux].
func NewServeMux() *ServeMux {
	return &ServeMux{
		mux: http.NewServeMux(),
	}
}
