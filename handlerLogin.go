package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type ResponseWithToken struct {
        ID uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email string `json:"email"`
        Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request){

        decoder := json.NewDecoder(r.Body)
        var params UserRequestParameters
        err := decoder.Decode(&params)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "Failed to verify user", err)
                return
        }

        user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user from database", err)
                return
        }

        err = auth.CheckPasswordHash([]byte(user.HashedPassword), params.Password)
        if err != nil{
                respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
                return
        }

        var expTime time.Duration
        if params.ExpiresInSeconds == nil || *params.ExpiresInSeconds > 3600 {
                expTime, err = time.ParseDuration("1h")
                if err != nil {
                        respondWithError(w, http.StatusInternalServerError, "Unable to set life duration for jwt", err)
                }
        } else {
                expTime = time.Duration(*params.ExpiresInSeconds)
        }

        jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expTime)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
        }

        responseStruct := ResponseWithToken{
                ID: user.ID,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Email: user.Email,
                Token: jwt,
        }

        respondWithJSON(w, http.StatusOK, responseStruct)
}
