package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/database"
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

	s := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error

	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
		chirps, err = cfg.dbQueries.GetAllChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching chirps")
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching chirps")
			return
		}
	}

	if r.URL.Query().Get("sort") == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
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
