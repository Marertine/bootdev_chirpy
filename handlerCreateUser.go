package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	//"github.com/lib/pq"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnSuccess struct {
		Id         string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Email      string `json:"email"`
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

	user, err := apiConfig.dbQueries.CreateUser(myCtx, myUserParams)
	if err != nil {
		// Type assertion to *pq.Error
		if pqErr, ok := err.(*pq.Error); ok {
			// Inspect the PostgreSQL error code
			fmt.Println("Postgres error code:", pqErr.Code)
			fmt.Println("Message:", pqErr.Message)
			fmt.Println("Detail:", pqErr.Detail)
			fmt.Println("Constraint:", pqErr.Constraint)

			// Example: unique violation
			if pqErr.Code == "23505" {
				return fmt.Errorf("User already exists")
			}
		}
		// All other errors
		return err
	}

	user, err := cfg.dbQueries.CreateUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Return the created user as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	respondWithJSON(w, 200, returnSuccess{Valid: true, Cleaned_body: strCleaned})

}
