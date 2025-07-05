package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/uller91/goChirpy/internal/aut"
	"github.com/uller91/goChirpy/internal/database"
)
	
type User struct {
	ID        	uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email     	string `json:"email"`
	Token	  	string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed    bool	`json:"is_chirpy_red"`
	//HashedPassword string `json: "hash"`	 for debugging only!
}


func (cfg *apiConfig) handlerChirpyRed(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
	}
	
	apiKey, err := aut.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error during token extraction", err)
		return
	} else if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong API key was sent", err)
		return
	}


	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {

		userID, err := uuid.Parse(params.Data.UserID)
		if err!= nil {
			respondWithError(w, http.StatusNotFound, "Non-uuid ID parsed", err)
			return
		} 

		_, err = cfg.database.UpdateCirpyRed(req.Context(), userID) 
		if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
    	}

		w.WriteHeader(http.StatusNoContent)
		return
	}
	
}


func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Email string `json:"email"`
		Password string  `json: "password"`
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

	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	returnUsr := User{}

	hashedPswd, err := aut.HashPassword(params.Password) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during password hashing", err)
		return
    }

	UsrParams := database.UpdateUserEmailPasswordParams{Email: params.Email, HashedPassword: hashedPswd, ID: userId}
	newUsr, err := cfg.database.UpdateUserEmailPassword(req.Context(), UsrParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the user creation", err)
		return
    }

	returnUsr.ID = newUsr.ID
	returnUsr.CreatedAt = newUsr.CreatedAt
	returnUsr.UpdatedAt = newUsr.UpdatedAt
	returnUsr.Email = newUsr.Email
	returnUsr.IsChirpyRed = newUsr.IsChirpyRed.Bool
	//returnUsr.HashedPassword = newUsr.HashedPassword	for debugging only!

	respondWithJSON(w, http.StatusOK, returnUsr)
}


func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Email string `json:"email"`
		Password string  `json: "password"`
    }

	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	returnUsr := User{}

	usr, err := cfg.database.GetUserFromEmail(req.Context(), params.Email)
	if err!= nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	} 

	err = aut.CheckPasswordHash(params.Password, usr.HashedPassword)
	if err!= nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	} 

	expire, _ := time.ParseDuration("1h")
	token, err := aut.MakeJWT(usr.ID, cfg.secret, expire) 
	if err!= nil {
		respondWithError(w, http.StatusInternalServerError, "Error during token creation", err)
		return
	} 

	expire, _ = time.ParseDuration("1440h")
	now := time.Now().UTC()
	expiresAt := now.Add(expire)
	refreshToken, _ := aut.MakeRefreshToken()

	refreshParams := database.CreateRefreshTokenParams{Token: refreshToken, UserID: usr.ID, ExpiresAt: expiresAt}
	rfrTkn, err := cfg.database.CreateRefreshToken(req.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the refresh token creation", err)
		return
    }

	returnUsr.ID = usr.ID
	returnUsr.CreatedAt = usr.CreatedAt
	returnUsr.UpdatedAt = usr.UpdatedAt
	returnUsr.Email = usr.Email
	returnUsr.Token = token
	returnUsr.RefreshToken = rfrTkn.Token
	returnUsr.IsChirpyRed = usr.IsChirpyRed.Bool

	respondWithJSON(w, http.StatusOK, returnUsr)
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
        Email string `json:"email"`
		Password string  `json: "password"`
    }

	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during JSON decoding", err)
		return
    }

	returnUsr := User{}

	hashedPswd, err := aut.HashPassword(params.Password) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during password hashing", err)
		return
    }

	UsrParams := database.CreateUserParams{Email: params.Email, HashedPassword: hashedPswd}
	newUsr, err := cfg.database.CreateUser(req.Context(), UsrParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was a problem during the user creation", err)
		return
    }

	returnUsr.ID = newUsr.ID
	returnUsr.CreatedAt = newUsr.CreatedAt
	returnUsr.UpdatedAt = newUsr.UpdatedAt
	returnUsr.Email = newUsr.Email
	returnUsr.IsChirpyRed = newUsr.IsChirpyRed.Bool
	//returnUsr.HashedPassword = newUsr.HashedPassword	for debugging only!

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