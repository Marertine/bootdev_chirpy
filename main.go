package main

import (
	"encoding/json"
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
		resetHandler(w, r, &apiCfg)
	})

	// Application handlers
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	// API handlers
	mux.HandleFunc("GET /api/healthz", healthzHandler)

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r, &apiCfg)
	})

	mux.HandleFunc("POST /api/validate_chirp", validationHandler)

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
	displayString := `<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(displayString, cfg.fileserverHits.Load())))
}

func resetHandler(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}

func validationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	/*type returnVals struct {
	    CreatedAt time.Time `json:"created_at"`
	    ID int `json:"id"`
	}*/

	type returnError struct {
		Error string `json:"error"`
	}

	type returnSuccess struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respDecodeErr := returnError{
			Error: "Something went wrong",
			//Error: "Invalid JSON in request body",
		}
		log.Printf("Error decoding parameters: %s", err)
		dat, _ := json.Marshal(respDecodeErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) == 0 {
		respBody := returnError{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		respBody := returnError{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	respBody := returnSuccess{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
