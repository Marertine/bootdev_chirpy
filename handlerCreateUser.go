package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnSuccess struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	// Parse the request body
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing JSON")
		return
	}

	// Create a new user in the database

	myCtx := context.Background()
	myUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     params.Email,
	}

	user, err := cfg.dbQueries.CreateUser(myCtx, myUserParams)
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

}
