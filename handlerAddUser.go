package main

import (
	"encoding/json"
	"net/http"
	"time"
        "github.com/TheGeneral00/Chirpy/internal/auth"
        "github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
        ID              uuid.UUID       `json:"id"`
        CreatedAt       time.Time       `json:"created_at"`
        UpdatedAt       time.Time       `json:"updated_at"`
        Email           string          `json:"email"`
        HashedPassword  string          `json:"hashed_password"`
}

type UserRequestParameters struct{
                Email                   string `json:"email"`
                Password                string `json:"password"`
        }

func (cfg *apiConfig) handlerAddUser (w http.ResponseWriter, r *http.Request){

        decoder := json.NewDecoder(r.Body)
        var params UserRequestParameters
        err := decoder.Decode(&params)
        if err != nil{
                respondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
                return
        }
        
        hashedPassword, err := auth.HashPassword(params.Password)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
                return
        }

        dbParams := database.CreateUserParams{
                Email: params.Email,
                HashedPassword: hashedPassword, 
        }
        user, err := cfg.dbQueries.CreateUser(r.Context(), dbParams)
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
                HashedPassword:       dbUser.HashedPassword,
        }
}
