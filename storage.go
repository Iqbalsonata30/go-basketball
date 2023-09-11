package main

import (
	"database/sql"
	"fmt"

	"github.com/iqbalsonata30/go-basketball/helper"
	_ "github.com/lib/pq"
)

type Storage interface {
	TeamStorage
}

type TeamStorage interface {
	CreateTeam(*Team) error
	GetAllTeams() ([]Team, error)
	GetTeamById(int) (*Team, error)
	DeleteTeam(int) error
	UpdateTeam(*Team, int) error
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
	if err := helper.Required(req.TeamName, req.Gender); err != nil {
		return err
	}
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

func (s *PostgresStore) GetTeamById(id int) (*Team, error) {
	query := "SELECT id,team_name,gender from teams where ID = $1;"
	result := s.db.QueryRow(query, id)
	var team Team
	if err := result.Scan(&team.ID, &team.TeamName, &team.Gender); err != nil {
		return nil, fmt.Errorf("team id's %d is not found ", id)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *PostgresStore) DeleteTeam(id int) error {
	query := "DELETE FROM teams where id = $1;"
	resp, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	res, err := resp.RowsAffected()
	if err != nil {
		return err
	}
	if res < 1 {
		return fmt.Errorf("team id's %d is not found ", id)
	}

	return nil
}

func (s *PostgresStore) UpdateTeam(req *Team, id int) error {
	query := "UPDATE teams SET team_name = $1,gender = $2 where ID = $3"
	resp, err := s.db.Exec(query, req.TeamName, req.Gender, id)
	if err != nil {
		return err
	}
	res, err := resp.RowsAffected()
	if err != nil {
		return err
	}
	if res < 1 {
		return fmt.Errorf("team id's %d is not found ", id)
	}
	return nil
}
