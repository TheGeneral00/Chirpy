package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
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

func (cfg *APIConfig) handlerLogin (w http.ResponseWriter, r *http.Request){
	
	requestID := r.Header.Clone().Get("X-request-ID")
	requestUUID, err := uuid.Parse(requestID)
	if err != nil{
		cfg.Logger.Failure.Printf("No valid X-request-ID retrieved: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve X-request-ID", err)
		return 
	}
        
	decoder := json.NewDecoder(r.Body)
        var params UserRequestParameters
        err = decoder.Decode(&params)
        if err != nil{
                helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to verify user", err)
                return
        }

        user, err := cfg.DBQueries.GetUserByEmail(r.Context(), params.Email)
        if err != nil{
                helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user from database", err)
		cfg.LogFailure(requestUUID, err)
                return
        }

        err = auth.CheckPasswordHash([]byte(user.HashedPassword), params.Password)
        if err != nil{
                helpers.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		cfg.LogFailure(requestUUID, err)
                return
        }

        jwtExpTime, err := time.ParseDuration("1h")
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
		cfg.LogFailure(requestUUID, err)
                return
        }
        jwt, err := auth.MakeJWT(user.ID, cfg.JWTSecret, jwtExpTime)
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
		cfg.LogFailure(requestUUID, err)
		log.Printf("Failed to create JWT token for Event: %d with error: %v", requestUUID, err)
                return
        }

        refreshTokenExpTime, err := time.ParseDuration("1440h")
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
		cfg.LogFailure(requestUUID, err)
                return
        }
        refreshToken, err := auth.MakeRefreshToken()
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
		cfg.LogFailure(requestUUID, err)
		log.Printf("Failed to create refresh token for Event: %d with error: %v", requestUUID, err)
                return
        }
        refreshTokenParams := database.CreateRefreshTokenParams{
                UserID: user.ID,
                Token: refreshToken,
                ExpiresAt: time.Now().Add(refreshTokenExpTime),
        }
        _, err = cfg.DBQueries.CreateRefreshToken(r.Context(), refreshTokenParams)

        responseStruct := ResponseWithToken{
                ID: user.ID,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Email: user.Email,
                Token: jwt,
                RefreshToken: refreshToken,
                IsChirpyRed: user.IsChirpyRed,
        }
	cfg.LogSuccess(requestUUID)
        helpers.RespondWithJSON(w, http.StatusOK, responseStruct)
}
