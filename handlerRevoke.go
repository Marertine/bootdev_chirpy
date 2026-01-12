package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	"github.com/Marertine/bootdev_chirpy/internal/database"
)

func handlerRevoke(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	// Use GetRefreshToken to have the needed parameters to revoke its access
	//refresh_token, err := auth.GetRefreshToken(r.Header)
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	//Use GetUserFromRefreshToken to have the needed parameters to revoke its access
	database_token, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		respondWithError(w, 500, "Something went wrong")
		return
	}
	/*
		if database_token.ExpiresAt.Before(time.Now().UTC()) {
			respondWithError(w, 401, "Unauthorized")
			return
		}*/

	// Don't care about validity, just revoke it
	myTokenParams := database.RevokeTokenParams{
		Token:     database_token.Token,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}

	err = cfg.dbQueries.RevokeToken(r.Context(), myTokenParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Return a 204 No Content status code
	respondWithJSON(w, 204, "")
}
