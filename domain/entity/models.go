package models

type Message struct {
	Username string `json:"username"`
	To       string `json:"to"`
	Message  string `json:"message"`
}
