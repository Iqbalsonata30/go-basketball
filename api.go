package main

type Teams interface {
	CreateTeam() error
	DeleteTeam(int) error
	GetAllTeams() ([]Teams, error)
}



