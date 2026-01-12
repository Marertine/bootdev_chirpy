package main

import (
	"database/sql"
	"net/http"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	//"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerDeleteChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	// Test with GetBearerToken to verify token still authorised
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	myidStr := r.URL.Path[len("/api/chirps/"):]
	myID, err := uuid.Parse(myidStr)
	if err != nil {
		respondWithError(w, 404, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), myID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Chirp not found")
			return
		}
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Check if the chirp belongs to the authenticated user
	if chirp.UserID != userID {
		respondWithError(w, 403, "Forbidden: You can only delete your own chirps")
		return
	}

	// Delete the chirp from the database
	err = cfg.dbQueries.DeleteChirpByID(r.Context(), myID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp")
		return
	}

	// Return success response (204: http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}
