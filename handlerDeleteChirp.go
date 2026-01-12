package main

import (
	"net/http"

	"github.com/Marertine/bootdev_chirpy/internal/auth"
	//"github.com/Marertine/bootdev_chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerDeleteChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	// Test with GetBearerToken to verify token still authorised
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 403, "Unauthorized")
		return
	}

	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 403, "Unauthorized")
		return
	}

	myidStr := r.URL.Path[len("/api/chirps/"):]
	myID, err := uuid.Parse(myidStr)
	if err != nil {
		respondWithError(w, 404, "Invalid chirp ID")
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
