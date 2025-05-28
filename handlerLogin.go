package main

import (
	"encoding/json"
	"net/http"

	"github.com/TheGeneral00/Chirpy/internal/auth"
)

func (cfg apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request){

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
        respondWithJSON(w, http.StatusOK, dbUserToUser(user))
}
