package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
)

type RefreshRequest struct{
        Token string `json:"token"`
}

func(cfg *APIConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){
        TokenString, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
        if TokenString == "" {
                helpers.RespondWithError(w, http.StatusUnauthorized, "No token string received", nil)
                return
        }
        token, err := cfg.dbQueries.RetrieveRefreshToken(r.Context(), TokenString)
        if err != nil {
                helpers.RespondWithError(w, http.StatusUnauthorized, "Refresh token unknown", err)
                return
        }
        if !token.ExpiresAt.After(time.Now()) || token.RevokedAt.Valid {
                helpers.RespondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked.", nil)
                return
        }
        expTime, err := time.ParseDuration("1h")
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create jwt token", err)
                return
        }
        newJWT, err := auth.MakeJWT(token.UserID, cfg.jwtSecret, expTime)
        helpers.RespondWithJSON(w, http.StatusOK, RefreshRequest{Token: newJWT})
}
