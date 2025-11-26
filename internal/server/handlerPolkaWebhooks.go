package server

import (
	"encoding/json"
	"net/http"

	"github.com/TheGeneral00/Chirpy/internal/auth"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

type webhookRequest struct{
        Event string `json:"event"`
        Data struct{
                UserID string `json:"user_id"` 
        } `json:"data"` 
}

func (cfg *APIConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request){
        ApiKey, err := auth.GetAPIKey(r.Header)
        if err != nil {
                helpers.RespondWithError(w, http.StatusUnauthorized, "Failed to retrive ApiKey", err)
                return
        }
        if ApiKey != cfg.polkaKey {
                helpers.RespondWithError(w, http.StatusUnauthorized, "API key mismatch", nil)
                return
        }

        decoder := json.NewDecoder(r.Body)
        var params webhookRequest
        err = decoder.Decode(&params)
        if err != nil{
                helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
                return
        }
        if params.Event != "user.upgraded"{
                helpers.RespondWithError(w, http.StatusNoContent, "Requested event not supported", err)
                return
        }
        userID := uuid.MustParse(params.Data.UserID)
        
        err = cfg.dbQueries.UpgradeToRed(r.Context(), userID)
        if err != nil {
                helpers.RespondWithError(w, http.StatusNotFound, "Failed to upgrade user", err)
                return
        }
        helpers.RespondWithJSON(w, http.StatusNoContent, nil)
}
