package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerCreateChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Body    string `json:"body"`
		User_id string `json:"user_id"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnSuccess struct {
		Valid        bool   `json:"valid"`
		Cleaned_body string `json:"cleaned_body"`
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

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Return the created user as JSON
	respondWithJSON(w, 201, returnSuccess{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})

	respondWithJSON(w, 200, returnSuccess{Valid: true, Cleaned_body: strCleaned})

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
