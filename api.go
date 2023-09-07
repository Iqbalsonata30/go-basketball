package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIServer struct {
	address string
	storage Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		address: listenAddr,
		storage: store,
	}
}

func (s *APIServer) Run() error {
	http.HandleFunc("/teams", HandleFunc(s.handleAPI))

	log.Printf("server is running on port %v\n", s.address)
	return http.ListenAndServe(s.address, nil)

}

func (s *APIServer) handleAPI(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.GetAllTeams(w, r)
	case "POST":
		return s.CreateTeam(w, r)
	default:
		return fmt.Errorf("method %s is not allowed", r.Method)
	}

}

func (s *APIServer) GetAllTeams(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	req := &TeamRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	err := s.storage.CreateTeam(req)
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusCreated, "team has been created succesfully")
}

func JSONEncode(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

func HandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
