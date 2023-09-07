package main

type TeamRequest struct {
	ID       int    `json:"id"`
	TeamName string `json:"teamName"`
	Gender   string `json:"gender"`
}
