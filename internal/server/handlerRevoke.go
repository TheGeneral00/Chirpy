package server

import (
	"net/http"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
)

func(cfg *APIConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
        TokenString, err := auth.GetBearerToken(r.Header)
        if err != nil {
                helpers.RespondWithError(w, http.StatusUnauthorized, "No token supplied", err)
                return
        }
        cfg.dbQueries.RevokeRefreshToken(r.Context(), TokenString)
        helpers.RespondWithJSON(w, http.StatusNoContent, nil)
}
