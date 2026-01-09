package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	//"github.com/google/uuid"
	//"github.com/Marertine/bootdev_chirpy/internal/config"
	//"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	fmt.Println("Boot.Dev/Twitter Clone")

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
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

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
