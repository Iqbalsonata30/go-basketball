package main

type Team struct {
	ID       int    `json:"id"`
	TeamName string `json:"teamName"`
	Gender   string `json:"gender"`
}
type ApiError struct {
	Error string `json:"error"`
}
