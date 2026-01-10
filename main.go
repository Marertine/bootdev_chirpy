package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
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

	platform := os.Getenv("PLATFORM")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
	}

	filepathRoot := "."
	mux := http.NewServeMux()

	// Admin handlers
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) { handlerMetrics(w, r, &apiCfg) })

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) { handlerReset(w, r, &apiCfg) })

	// Application handlers
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	// API handlers
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) { handlerGetAllChirps(w, r, &apiCfg) })

	mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) { handlerGetChirpByID(w, r, &apiCfg) })

	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) { handlerCreateChirp(w, r, &apiCfg) })

	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) { handlerLogin(w, r, &apiCfg) })

	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) { handlerCreateUser(w, r, &apiCfg) })

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
