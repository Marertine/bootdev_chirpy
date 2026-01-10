package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	"github.com/google/uuid"
)

func handlerLogin(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnSuccess struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if match != true {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	// Return the user as JSON
	respondWithJSON(w, 200, returnSuccess{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})

}
