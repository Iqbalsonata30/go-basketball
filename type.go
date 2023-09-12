package main

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type Team struct {
	ID       int    `json:"id"`
	TeamName string `json:"teamName"`
	Gender   string `json:"gender"`
}

type Player struct {
	ID        uuid.UUID  `json:"id"`
	TeamID    int        `json:"teamID"`
	Name      string     `json:"name"`
	Number    int        `json:"number"`
	Height    int        `json:"height"`
	Birthdate civil.Date `json:"birthdate"`
}

type ApiError struct {
	Error string `json:"error"`
}
