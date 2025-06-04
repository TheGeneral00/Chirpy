package main

import (
	"net/http"
	"strings"
)

func(cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
        TokenString, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
        if TokenString == "" {
                respondWithError(w, http.StatusInternalServerError, "No token string received", nil)
                return
        }
        cfg.dbQueries.RevokeRefreshToken(r.Context(), TokenString)
        respondWithJSON(w, http.StatusNoContent, nil)
}
