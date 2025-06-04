package main

import (
	"net/http"
	"github.com/TheGeneral00/Chirpy/internal/auth"
)

func(cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
        TokenString, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "No token supplied", err)
                return
        }
        cfg.dbQueries.RevokeRefreshToken(r.Context(), TokenString)
        respondWithJSON(w, http.StatusNoContent, nil)
}
