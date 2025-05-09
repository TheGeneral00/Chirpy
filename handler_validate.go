package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
                CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
                CleanedBody: handleProfanity(params.Body),
	})
}

//This function takes the body of a chirp and scans for profanity. If it finds some it is replaced by four asterix
func handleProfanity (body string) string {
        profList := map[string]bool{
                "kerfuffle": true,
                "sharbert": true,
                "fornax": true,
        } 
        words := strings.Split(body, " ")

        for i, word := range words {
                lowerWord := strings.ToLower(word)
                _, isProfane := profList[lowerWord]
                if isProfane {
                        words[i] = "****"
                }
        }
        return strings.Join(words, " ")
}
