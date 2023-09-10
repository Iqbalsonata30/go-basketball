package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/teams", HandleFunc(s.handleAPI))
	router.HandleFunc("/teams/{id}", HandleFunc(s.handleAPIById))
	log.Printf("server is running on port %v\n", s.address)

	http.ListenAndServe(s.address, router)
}

func (s *APIServer) handleAPI(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.GetAllTeams(w, r)
	}
	if r.Method == "POST" {
		return s.CreateTeam(w, r)
	}
	return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
}

func (s *APIServer) handleAPIById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	teamID, _ := strconv.Atoi(id)
	if r.Method == "GET" {
		return s.GetTeamById(w, r, teamID)
	}
	return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
}

func (s *APIServer) GetAllTeams(w http.ResponseWriter, r *http.Request) error {
	teams, err := s.storage.GetAllTeams()
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusOK, teams)

}

func (s *APIServer) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	req := new(Team)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	err := s.storage.CreateTeam(req)
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusCreated, "team has been created succesfully")
}

func (s *APIServer) GetTeamById(w http.ResponseWriter, r *http.Request, id int) error {
	team, err := s.storage.GetTeamById(id)
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusOK, team)
}

func JSONEncode(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func HandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			JSONEncode(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
