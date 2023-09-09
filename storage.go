package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateTeam(*Team) error
	GetAllTeams() ([]Team, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://postgres:secret@localhost:5433/go_basketball?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	query :=
		`CREATE TABLE IF NOT EXISTS teams(
		id SERIAL PRIMARY KEY,
		team_name VARCHAR(100) NOT NULL,
		gender VARCHAR(50) NOT NULL
	);`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateTeam(req *Team) error {
	query := "INSERT INTO teams(team_name,gender) VALUES($1,$2);"
	_, err := s.db.Exec(query, req.TeamName, req.Gender)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStore) GetAllTeams() ([]Team, error) {
	query := "SELECT id,team_name,gender from teams;"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team

	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.ID, &team.TeamName, &team.Gender); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
