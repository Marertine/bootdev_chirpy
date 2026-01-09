package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerCreateChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) == 0 {
		respondWithError(w, 400, "Chirp body is required")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	strCleaned := sanitizeString(params.Body)

	// Valid, insert into DB
	myChirpParams := database.CreateChirpParams{
		Body:   strCleaned,
		UserID: params.UserID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), myChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	// Return the created chirp as JSON
	respondWithJSON(w, 201, returnSuccess{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func sanitizeString(input string) string {
	cleanedWords := strings.Split(input, " ")
	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

	for i, word := range cleanedWords {
		if _, found := badWords[strings.ToLower(word)]; found {
			cleanedWords[i] = "****"
		}
	}

	strCleaned := strings.Join(cleanedWords, " ")

	return strCleaned

}
