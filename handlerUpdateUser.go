package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerUpdateUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnSuccess struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	// Test with GetRefreshToken to verify token still authorised
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	//GetUserFromRefreshToken to obtain user identity
	database_token, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Parse the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	params.Password = hashedPassword

	// Valid, update user details in the database
	myUpdateParams := database.UpdateUserParams{
		ID:             database_token.ID,
		HashedPassword: params.Password,
		Email:          params.Email,
		UpdatedAt:      time.Now().UTC(),
	}

	// Update the user in the database
	user, err := cfg.dbQueries.UpdateUser(r.Context(), myUpdateParams)
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
