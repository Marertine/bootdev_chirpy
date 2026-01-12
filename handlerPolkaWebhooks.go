package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_ID uuid.UUID `json:"user_id"`
		}
	}

	// Test with GetAPIKey to verify APIkey is correct
	request_apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	if request_apikey != cfg.polkaAPIKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// Parse the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		// Return a 204 No Content status code
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Update the user in the database
	_, err = cfg.dbQueries.UpdateIsRed(r.Context(), database.UpdateIsRedParams{
		ID:          params.Data.User_ID,
		IsChirpyRed: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "User not found")
			return
		}
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Return a 204 No Content status code to indicate success
	w.WriteHeader(http.StatusNoContent)

}
