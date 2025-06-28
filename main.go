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
)


type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
	platform string
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


func main() {
	godotenv.Load()	//to load environmental variables
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	var counter atomic.Int32
	counter.Store(0)
	apiCfg := &apiConfig{fileserverHits: counter, database: dbQueries, platform: os.Getenv("PLATFORM")}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerOk)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerReq)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	
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