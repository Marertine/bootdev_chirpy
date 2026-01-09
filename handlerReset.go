package main

import (
	"net/http"
)

func handlerReset(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {

	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error resetting database")
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset and users deleted"))
}
