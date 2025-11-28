package server

import (
	"net/http"

	"github.com/TheGeneral00/Chirpy/internal/helpers"
)

func (cfg *APIConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
        err := cfg.DBQueries.ResetUsers(r.Context())
        if err != nil {
                helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to reset database", err)
                return
        }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

