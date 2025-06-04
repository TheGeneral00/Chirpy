package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
        ID              uuid.UUID       `json:"id"`
        CreatedAt       time.Time       `json:"created_at"`
        UpdatedAt       time.Time       `json:"updated_at"`
        Body            string          `json:"body"`
        UserID          uuid.UUID       `json:"user_id"`
}

type dbParams struct {
        Body            string          `json:"body"`
        UserID          uuid.UUID       `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to retrieve token", err)
        }
        userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Failed to validate token", err)
        }

        if userID == uuid.Nil {
                respondWithError(w, http.StatusUnauthorized, "Token is not valid. Access denied.", nil)
        }
        
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

        dbParams := database.CreateChirpParams{
                Body: params.Body,
                UserID: userID,
        }
        
        dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), dbParams)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp in db", err)
                return 
        }
        chirp := dbChirpToChirp(dbChirp)
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
    dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Unable to get all chirps", err)
        return
    }
    
    // Convert database chirps to API chirps with correct JSON field names
    chirps := []Chirp{}
    for _, dbChirp := range dbChirps {
        chirps = append(chirps, dbChirpToChirp(dbChirp))
    }
    
    respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
        
        idString := r.PathValue("chirpID")
        if idString == "" {
                respondWithJSON(w, http.StatusOK, Chirp{})
                return
        }

        id, err := uuid.Parse(idString)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "unable to parse id", err)
        }

        dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), id)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "unable to retrieve chirp", err)
                return
        }

        respondWithJSON(w, http.StatusOK, dbChirpToChirp(dbChirp))
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

func dbChirpToChirp ( dbChirp database.Chirp) Chirp {
        return Chirp {
                ID: dbChirp.ID,
                CreatedAt: dbChirp.CreatedAt,
                UpdatedAt: dbChirp.UpdatedAt,
                Body: dbChirp.Body,
                UserID: dbChirp.UserID,
        }
}
