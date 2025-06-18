package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

type ResponseWithToken struct {
        ID uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email string `json:"email"`
        Token string `json:"token"`
        RefreshToken string `json:"refresh_token"`
        IsChirpyRed bool `json:"is_chirpy_red"`
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

        jwtExpTime, err := time.ParseDuration("1h")
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
                return
        }
        jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, jwtExpTime)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
                return
        }

        refreshTokenExpTime, err := time.ParseDuration("1440h")
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
                return
        }
        refreshToken, err := auth.MakeRefreshToken()
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
                return
        }
        refreshTokenParams := database.CreateRefreshTokenParams{
                UserID: user.ID,
                Token: refreshToken,
                ExpiresAt: time.Now().Add(refreshTokenExpTime),
        }
        _, err = cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)

        responseStruct := ResponseWithToken{
                ID: user.ID,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Email: user.Email,
                Token: jwt,
                RefreshToken: refreshToken,
                IsChirpyRed: user.IsChirpyRed,
        }

        respondWithJSON(w, http.StatusOK, responseStruct)
}
