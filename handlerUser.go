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
        IsChirpyRed     bool            `json:"is_chirpy_red"`
}

type UserRequestParameters struct{
                Email                   string `json:"email"`
                Password                string `json:"password"`
}

type EmailResponse struct{
        Email string `json:"email"`
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
                IsChirpyRed:    dbUser.IsChirpyRed,
        }
}

func(cfg *apiConfig) handlerUpdateUserCredentials(w http.ResponseWriter, r *http.Request) {
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "No Token retrieved", err)
                return
        }
        userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Token not valid", err)
                return
        }

        decoder := json.NewDecoder(r.Body)
        var params UserRequestParameters
        err = decoder.Decode(&params)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to read request body", err)
                return
        }

        hashedPassword, err := auth.HashPassword(params.Password)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to hash Password", err)
                return
        }

        newUserParams := database.UpdateUserCredentialsParams{
                ID: userID,
                Email: params.Email,
                HashedPassword: hashedPassword,
        }
        
        newUserCredentials, err := cfg.dbQueries.UpdateUserCredentials(r.Context(), newUserParams)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to update user credentials", err)
                return
        }
        respondWithJSON(w, http.StatusOK, EmailResponse{Email: newUserCredentials.Email})
}
