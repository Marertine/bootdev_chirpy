package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerGetAllChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
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

	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps")
		return
	}

	var response []returnSuccess
	for _, chirp := range chirps {
		response = append(response, returnSuccess{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, response)
}
