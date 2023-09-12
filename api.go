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

	router.HandleFunc("/teams", HandleFunc(s.handleTeamAPI))
	router.HandleFunc("/teams/{id}", HandleFunc(s.handleTeamAPIById))
	router.HandleFunc("/players", HandleFunc(s.handlePlayerAPI))

	log.Printf("server is running on port %v\n", s.address)
	http.ListenAndServe(s.address, router)
}

func (s *APIServer) handleTeamAPI(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.FindAllTeams(w, r)
	}
	if r.Method == "POST" {
		return s.CreateTeam(w, r)
	}
	return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
}

func (s *APIServer) handleTeamAPIById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	teamID, _ := strconv.Atoi(id)
	if r.Method == "GET" {
		return s.FindTeamById(w, r, teamID)
	}
	if r.Method == "DELETE" {
		return s.DeleteTeam(w, r, teamID)
	}
	if r.Method == "PUT" {
		return s.UpdateTeam(w, r, teamID)
	}
	return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
}

func (s *APIServer) handlePlayerAPI(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.CreatePlayer(w, r)
	default:
		return JSONEncode(w, http.StatusBadRequest, ApiError{Error: "method " + r.Method + " is not allowed"})
	}
}

func (s *APIServer) FindAllTeams(w http.ResponseWriter, r *http.Request) error {
	teams, err := s.storage.FindAllTeams()
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
	if err := s.storage.CreateTeam(req); err != nil {
		return err
	}

	return JSONEncode(w, http.StatusCreated, helper.WriteMessageAPI(http.StatusCreated, "Team has been created sucesfully."))
}

func (s *APIServer) FindTeamById(w http.ResponseWriter, r *http.Request, id int) error {
	team, err := s.storage.FindTeamById(id)
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusOK, team)
}

func (s *APIServer) DeleteTeam(w http.ResponseWriter, r *http.Request, id int) error {
	err := s.storage.DeleteTeam(id)
	if err != nil {
		return err
	}
	return JSONEncode(w, http.StatusOK, helper.WriteMessageAPI(http.StatusOK, "Team has been deleted sucesfully.he"))
}

func (s *APIServer) UpdateTeam(w http.ResponseWriter, r *http.Request, id int) error {
	req := new(Team)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	if err := s.storage.UpdateTeam(req, id); err != nil {
		return err
	}
	return JSONEncode(w, http.StatusOK, helper.WriteMessageAPI(http.StatusOK, "Team has been updated sucesfully."))
}

func (s *APIServer) CreatePlayer(w http.ResponseWriter, r *http.Request) error {
	req := new(Player)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	if err := s.storage.CreatePlayer(req); err != nil {
		return err
	}
	return JSONEncode(w, http.StatusCreated, helper.WriteMessageAPI(http.StatusCreated, "Player has been created succesfully."))

}

func JSONEncode(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

func HandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			JSONEncode(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
