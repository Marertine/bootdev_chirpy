package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	/*type returnVals struct {
	    CreatedAt time.Time `json:"created_at"`
	    ID int `json:"id"`
	}*/

	type returnError struct {
		Error string `json:"error"`
	}

	type returnSuccess struct {
		Valid        bool   `json:"valid"`
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		/*
			respDecodeErr := returnError{
				Error: "Something went wrong",
				//Error: "Invalid JSON in request body",
			}
			log.Printf("Error decoding parameters: %s", err)
			dat, _ := json.Marshal(respDecodeErr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write(dat)
		*/
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) == 0 {
		/*
			respBody := returnError{
				Error: "Something went wrong",
			}
			dat, err := json.Marshal(respBody)
			if err != nil {
				log.Printf("Error marshalling JSON: %s", err)
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(dat)
		*/
		respondWithError(w, 400, "Chirp body is required")
		return
	}

	if len(params.Body) > 140 {
		/*
			respBody := returnError{
				Error: "Chirp is too long",
			}
			dat, err := json.Marshal(respBody)
			if err != nil {
				log.Printf("Error marshalling JSON: %s", err)
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(dat)
		*/
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	/*
		respBody := returnSuccess{
			Valid: true,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	*/

	strCleaned := sanitizeString(params.Body)
	respondWithJSON(w, 200, returnSuccess{Valid: true, Cleaned_body: strCleaned})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func sanitizeString(input string) string {

	cleanedWords := strings.Split(input, " ")
	//badWords := []string{"kerfuffle", "sharbert", "fornax"}
	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

	for i, word := range cleanedWords {
		if _, found := badWords[strings.ToLower(word)]; found {
			cleanedWords[i] = "****"
		}

		/*
			for _, badWord := range badWords {
				if strings.EqualFold(word, badWord) {
					//replacement := strings.Repeat("*", len(badWord))
					cleanedWords[i] = "****"
				}
			}
		*/
	}

	strCleaned := strings.Join(cleanedWords, " ")

	return strCleaned

}
