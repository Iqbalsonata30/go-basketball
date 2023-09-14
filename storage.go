package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/iqbalsonata30/go-basketball/helper"
	_ "github.com/lib/pq"
)

type Storage interface {
	TeamStorage
	PlayerStorage
}

type TeamStorage interface {
	CreateTeam(*Team) error
	FindAllTeams() ([]Team, error)
	FindTeamById(int) (*Team, error)
	DeleteTeam(int) error
	UpdateTeam(*Team, int) error
}

type PlayerStorage interface {
	CreatePlayer(*Player) error
	FindAllPlayers() ([]Player, error)
    FindPlayerById(string) (*Player,error)
    DeletePlayer(string) error
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
	if err := s.createTableTeams(); err != nil {
		return err
	}
	if err := s.createTablePlayers(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) createTableTeams() error {
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

func (s *PostgresStore) createTablePlayers() error {
	query := `CREATE TABLE IF NOT EXISTS players(
		id uuid PRIMARY KEY,
		team_id integer REFERENCES teams(id) ON DELETE CASCADE,
		name varchar(150) NOT NULL,
		number integer NOT NULL,
		height integer NOT NULL
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

func (s *PostgresStore) FindAllTeams() ([]Team, error) {
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

func (s *PostgresStore) FindTeamById(id int) (*Team, error) {
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

func (s *PostgresStore) CreatePlayer(req *Player) error {
	if err := helper.Required(req.TeamID, req.Name, req.Number, req.Height); err != nil {
		return err
	}
	query := "INSERT INTO players(id,team_id,name,number,height) VALUES ($1,$2,$3,$4,$5);"
	_, err := s.db.Exec(query, uuid.New(), req.TeamID, req.Name, req.Number, req.Height)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) FindAllPlayers() ([]Player, error) {
	query := "SELECT id,team_id,name,number,height from players"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		err := rows.Scan(&p.ID, &p.TeamID, &p.Name, &p.Number, &p.Height)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return players, err
}
func (s *PostgresStore) DeletePlayer(id string) error{
    pID,err := uuid.Parse(id)
    if err != nil{
         return fmt.Errorf("player's id is not found")
    }
    query := "DELETE FROM players WHERE id = $1";
    _,err = s.db.Exec(query,pID)
    if err != nil{
        return err
    }
    return nil  
}

func (s *PostgresStore) FindPlayerById(id string) (*Player,error){
    pID,err := uuid.Parse(id)
    if err != nil{
        return nil,fmt.Errorf("player's id is not found") 
    }
    query := "SELECT id,team_id,name,number,height FROM players WHERE id = $1;"
    row := s.db.QueryRow(query,pID)
    var p Player 
    if err := row.Scan(&p.ID,&p.TeamID,&p.Name,&p.Number,&p.Height);err != nil{
        return nil,err
    }
    if err := row.Err();err != nil{
        return nil,err
    }
    return &p,nil
}
