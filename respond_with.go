package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func respondWithError(w http.ResponseWriter, respCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	log.Printf("Error happened: %s\n", msg)

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, respCode, errorResponse{Error: msg,})
}

func respondWithJSON(w http.ResponseWriter, respCode int, respBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	

	dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
	}

    w.WriteHeader(respCode)
    w.Write(dat)
}