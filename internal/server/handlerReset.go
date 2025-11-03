package server

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
        err := cfg.dbQueries.ResetUsers(r.Context())
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to reset database", err)
                return
        }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

