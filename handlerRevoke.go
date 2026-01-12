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
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Don't care about validity, just revoke it
	err = cfg.dbQueries.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     refresh_token,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	// Return a 204 No Content status code
	w.WriteHeader(http.StatusNoContent)
}
