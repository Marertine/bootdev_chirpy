package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
)

func handlerRefresh(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type returnSuccess struct {
		Token string `json:"token"`
	}

	// Test with GetBearerToken to verify user identity
	refresh_token, err := auth.GetRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	//GetUserFromRefreshToken to verify user identity
	database_token, err := auth.GetUserFromRefreshToken(r.Context(), refresh_token, cfg.dbQueries)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if database_token.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	token, err := auth.MakeJWT(database_token.UserID, cfg.secret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Return the token as JSON
	respondWithJSON(w, 200, returnSuccess{
		Token: token,
	})
}
