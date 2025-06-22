package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)
	
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Email string `json:"email"`
    }

	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	returnUsr := User{}
	newUsr, err := cfg.database.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the user creation", err)
		return
    }

	returnUsr.ID = newUsr.ID
	returnUsr.CreatedAt = newUsr.CreatedAt
	returnUsr.UpdatedAt = newUsr.UpdatedAt
	returnUsr.Email = newUsr.Email

	respondWithJSON(w, http.StatusCreated, returnUsr)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	
	cfg.fileserverHits.Store(0)

	err := cfg.database.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the users deletion", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}