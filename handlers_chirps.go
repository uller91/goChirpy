package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/uller91/goChirpy/internal/database"
)

type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}


func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := (req.PathValue("chirpID"))
	chirpUuid, err := uuid.Parse(chirpId)
	if err!= nil {
		respondWithError(w, http.StatusNotFound, "Non-uuid ID parsed", err)
		return
	} 

	chirp, err :=  cfg.database.GetChirp(req.Context(), chirpUuid)
	if err!= nil {
		respondWithError(w, http.StatusNotFound, "Chirp with this ID has not been found", err)
		return
	} 

	returnChirp := Chirp{}
	returnChirp.ID = chirp.ID
	returnChirp.CreatedAt = chirp.CreatedAt
	returnChirp.UpdatedAt = chirp.UpdatedAt
	returnChirp.Body = chirp.Body
	returnChirp.UserID = chirp.UserID

	respondWithJSON(w, http.StatusOK, returnChirp)
}


func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.database.GetChirps(req.Context())
	if err!= nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the chirps retrieval", err)
		return
	} 

	returnChirps := []Chirp{}
	chirpChirp := Chirp{}
	for _, chirp := range chirps {
		chirpChirp.ID = chirp.ID
		chirpChirp.CreatedAt = chirp.CreatedAt
		chirpChirp.UpdatedAt = chirp.UpdatedAt
		chirpChirp.Body = chirp.Body
		chirpChirp.UserID = chirp.UserID

		returnChirps = append(returnChirps, chirpChirp)
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Message string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	if len(params.Message) >= 140 {
		respondWithError(w, http.StatusBadRequest, "The message is too long!", err)
		return
	}

	//cleaning of the string
	cleanedMessage := profanityCheck(params.Message)

	returnChirp := Chirp{}

	dbParams := database.CreateChirpParams{Body: cleanedMessage, UserID: params.UserId}
	newChirp, err := cfg.database.CreateChirp(req.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the chirp creation", err)
		return
    }

	returnChirp.ID = newChirp.ID
	returnChirp.CreatedAt = newChirp.CreatedAt
	returnChirp.UpdatedAt = newChirp.UpdatedAt
	returnChirp.Body = newChirp.Body
	returnChirp.UserID = newChirp.UserID

	respondWithJSON(w, http.StatusCreated, returnChirp)
}

func profanityCheck(s string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(s, " ")
	for i, word := range words {
		for _, badWord := range profanity {
			if strings.ToLower(word) == badWord {
				words[i] = "****"
			}
		}
		
	}

	cleanMessage := strings.Join(words, " ")
	return cleanMessage
}