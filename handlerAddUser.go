package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
        ID              uuid.UUID       `json:"id"`
        CreatedAt       time.Time       `json:"created_at"`
        UpdatedAt       time.Time       `json:"updated_at"`
        Email           string          `json:"email"`
}

func (cfg *apiConfig) handlerAddUser (w http.ResponseWriter, r *http.Request){
        type parameters struct{
                Email string `json:"email"`
        }
        decoder := json.NewDecoder(r.Body)
        params := parameters{}
        err := decoder.Decode(&params)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
                return
        }
        user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
                return
        }
        respondWithJSON(w, http.StatusCreated, dbUserToUser(user))
}

func dbUserToUser (dbUser database.User) User {
        return User {
                ID:             dbUser.ID,
                CreatedAt:      dbUser.CreatedAt,
                UpdatedAt:      dbUser.UpdatedAt,
                Email:          dbUser.Email,
        }
}
