package web_socket_service

import (
	"log"
	"net/http"
	"sd-chat/domain/entity"

	"github.com/gorilla/websocket"
)

type ServerConnData struct {
	Clients   map[string]*websocket.Conn // username → connection
	Broadcast chan models.Message
}

func NewServerConnData() *ServerConnData {
	return &ServerConnData{
		Clients:   make(map[string]*websocket.Conn),
		Broadcast: make(chan models.Message),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections upgrades HTTP requests to WebSocket and registers clients by username
func (scd *ServerConnData) HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("[WebSocket] Attempting to upgrade connection...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()
	log.Println("[WebSocket] Connection upgraded.")

	var intro models.Message
	err = conn.ReadJSON(&intro)
	if err != nil || intro.Username == "" {
		log.Printf("[WebSocket] Invalid initial message or missing username: %v\n", err)
		return
	}

	username := intro.Username
	log.Printf("[WebSocket] User connected: %s\n", username)

	// Register user connection
	scd.Clients[username] = conn

	defer func() {
		log.Printf("[WebSocket] User disconnected: %s\n", username)
		delete(scd.Clients, username)
	}()

	// Message read loop
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("[WebSocket] Error reading message from %s: %v\n", username, err)
			return
		}

		log.Printf("[WebSocket] Received message from %s to %s: %s\n", msg.Username, msg.To, msg.Message)
		scd.Broadcast <- msg
	}
}

// HandleMessages sends messages to the intended user or broadcasts if msg.To == "all"
func (scd *ServerConnData) HandleMessages() {
	for {
		msg := <-scd.Broadcast

		if msg.To != "" && msg.To != "all" {
			// Send to specific user
			conn, ok := scd.Clients[msg.To]
			if ok {
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Printf("[WebSocket] Error sending message to %s: %v\n", msg.To, err)
					conn.Close()
					delete(scd.Clients, msg.To)
				} else {
					log.Printf("[WebSocket] Sent message to %s\n", msg.To)
					// Send it back to your own chat for feedback
					conn = scd.Clients[msg.Username]
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("[WebSocket] Error sending own feedback message to %s: %v\n", msg.To, err)
						conn.Close()
						delete(scd.Clients, msg.To)
					}
				}
			} else {
				log.Printf("[WebSocket] Tried to send message to unknown user: %s\n", msg.To)
			}
		} else {
			// Broadcast to all users
			for username, conn := range scd.Clients {
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Printf("[WebSocket] Error broadcasting to %s: %v\n", username, err)
					conn.Close()
					delete(scd.Clients, username)
				}
			}
			log.Printf("[WebSocket] Broadcasted message from %s\n", msg.Username)
		}
	}
}
