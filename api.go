package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iqbalsonata30/go-basketball/helper"
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
	if r.Method == "DELETE" {
		return s.DeleteTeam(w, r, teamID)
	}
	return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
}

func (s *APIServer) GetAllTeams(w http.ResponseWriter, r *http.Request) error {
	teams, err := s.storage.GetAllTeams()
	helper.CheckError(err)
	return JSONEncode(w, http.StatusOK, teams)

}

func (s *APIServer) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	req := new(Team)

	err := json.NewDecoder(r.Body).Decode(req)
	helper.CheckError(err)

	err = s.storage.CreateTeam(req)
	helper.CheckError(err)

	return JSONEncode(w, http.StatusCreated, helper.WriteMessageAPI(http.StatusCreated, "Team has been created sucesfully."))
}

func (s *APIServer) GetTeamById(w http.ResponseWriter, r *http.Request, id int) error {
	team, err := s.storage.GetTeamById(id)
	helper.CheckError(err)
	return JSONEncode(w, http.StatusOK, team)
}

func (s *APIServer) DeleteTeam(w http.ResponseWriter, r *http.Request, id int) error {
	err := s.storage.DeleteTeam(id)
	helper.CheckError(err)
	return JSONEncode(w, http.StatusOK, helper.WriteMessageAPI(http.StatusOK, "Team has been deleted sucesfully.he"))
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
