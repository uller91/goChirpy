package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/uller91/goChirpy/internal/database"
	"github.com/uller91/goChirpy/internal/aut"
	"sort"
)

type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}


func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := (req.PathValue("chirpID"))
	chirpUuid, err := uuid.Parse(chirpId)
	if err!= nil {
		respondWithError(w, http.StatusNotFound, "Non-uuid ID parsed", err)
		return
	} 

	token, err := aut.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error during token extraction", err)
		return
	}

	userId, err := aut.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token presented", err)
		return
	}

	chirp, err :=  cfg.database.GetChirp(req.Context(), chirpUuid)
	if err!= nil {
		respondWithError(w, http.StatusNotFound, "Chirp with this ID has not been found", err)
		return
	} 

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Trying to delete the Chirp which doesn't belong to you", err)
	} else {
		
	err = cfg.database.DeleteChirp(req.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the chirp deletion", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	}
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
	authorId := req.URL.Query().Get("author_id")

	chirps := []database.Chirp{}
	returnChirps := []Chirp{}
	chirpChirp := Chirp{}
	var err error

	if authorId != "" {
		authorUuid, err := uuid.Parse(authorId)
		if err!= nil {
			respondWithError(w, http.StatusNotFound, "Non-uuid ID parsed", err)
			return
		} 

		chirps, err = cfg.database.GetChirpsByAuthor(req.Context(), authorUuid)
		if err!= nil {
			respondWithError(w, http.StatusInternalServerError, "There was a problem during the chirps retrieval", err)
			return
		} 
	} else {
		chirps, err = cfg.database.GetChirps(req.Context())
		if err!= nil {
			respondWithError(w, http.StatusInternalServerError, "There was a problem during the chirps retrieval", err)
			return
		} 
	}
	
	for _, chirp := range chirps {
		chirpChirp.ID = chirp.ID
		chirpChirp.CreatedAt = chirp.CreatedAt
		chirpChirp.UpdatedAt = chirp.UpdatedAt
		chirpChirp.Body = chirp.Body
		chirpChirp.UserID = chirp.UserID

		returnChirps = append(returnChirps, chirpChirp)
	}

	order := req.URL.Query().Get("sort")
	if order == "desc" {
		sort.Slice(returnChirps, func(i, j int) bool {return returnChirps[j].CreatedAt.Before(returnChirps[i].CreatedAt)})
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Message string `json:"body"`
		UserId uuid.UUID `json:"user_id"` //unused at the moment
    }

	token, err := aut.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No autentication token in the header", err)
		return
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	userId, err := aut.ValidateJWT(token, cfg.secret)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "No access for you!", err)
		return
    }

	if len(params.Message) >= 140 {
		respondWithError(w, http.StatusBadRequest, "The message is too long!", err)
		return
	}

	//cleaning of the string
	cleanedMessage := profanityCheck(params.Message)

	returnChirp := Chirp{}

	dbParams := database.CreateChirpParams{Body: cleanedMessage, UserID: userId}
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