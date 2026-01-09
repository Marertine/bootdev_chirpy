package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerGetChirpByID(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type returnError struct {
		Error string `json:"error"`
	}

	type returnSuccess struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	myidStr := r.URL.Path[len("/api/chirps/"):]
	myID, err := uuid.Parse(myidStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), myID)
	if err != nil {
		respondWithError(w, 404, "Error fetching chirp")
		return
	}

	response := returnSuccess{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 200, response)
}
