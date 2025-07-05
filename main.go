package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"github.com/uller91/goChirpy/internal/database"
	"github.com/joho/godotenv"
	"time"
	"github.com/uller91/goChirpy/internal/aut"
)


type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
	platform string
	secret string
	polkaKey string
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) handlerReq(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
}


func handlerOk(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type Reply struct {
        Token string `json:"token"`
    }

	RefreshToken, err := aut.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error during token extraction", err)
		return
	}

	usrTkn, err := cfg.database.GetUserFromRefreshToken(req.Context(), RefreshToken)
	if err!= nil || time.Now().UTC().After(usrTkn.TokenExpiresAt) || usrTkn.TokenRevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect or expired token", err)
		return
	} 

	expire, _ := time.ParseDuration("1h")
	token, err := aut.MakeJWT(usrTkn.ID, cfg.secret, expire) 
	if err!= nil {
		respondWithError(w, http.StatusInternalServerError, "Error during token creation", err)
		return
	} 

	reply := Reply{Token: token,}
	respondWithJSON(w, http.StatusOK, reply)
}


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	RefreshToken, err := aut.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during token extraction", err)
		return
	}

	_, err = cfg.database.RevokeRefreshToken(req.Context(), RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error during token revoking", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func main() {
	godotenv.Load()	//to load environmental variables
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	var counter atomic.Int32
	counter.Store(0)
	apiCfg := &apiConfig{fileserverHits: counter, database: dbQueries, platform: os.Getenv("PLATFORM"), secret: os.Getenv("SECRET"), polkaKey: os.Getenv("POLKA_KEY")}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerOk)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerReq)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerChirpyRed)

	s := http.Server{
		Addr: ":8080", 
		Handler: mux,
	}

	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

	return
}