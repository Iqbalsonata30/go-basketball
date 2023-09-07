package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIServer struct {
	address string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		address: listenAddr,
	}
}

func (s *APIServer) Run() error {
	http.HandleFunc("/teams", HandleFunc(s.handleAPI))

	log.Printf("server is running on port %v\n", s.address)
	return http.ListenAndServe(s.address, nil)

}

func (s *APIServer) handleAPI(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.GetAllTeams(w, r)
	}
	return fmt.Errorf("method %s is not allowed", r.Method)
}

func (s *APIServer) GetAllTeams(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode("Succeddd")
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

func HandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
