package web_socket_service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sd-chat/domain/entity"
)

type ServerConnData struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan models.Message
}

func NewServerConnData() *ServerConnData {
	return &ServerConnData{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan models.Message),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (scd *ServerConnData) HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Opening a new connection...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection | Error: %v", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection upgraded successfully!")

	scd.Clients[conn] = true

	fmt.Printf("Waiting for messages on connection...\n")
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading message | Error: %v\n", err)
			delete(scd.Clients, conn)
			return
		}

		scd.Broadcast <- msg
	}
}

func (scn *ServerConnData) HandleMessages() {
	for {
		msg := <-scn.Broadcast

		for client := range scn.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Error writing message to client | err: %v\n", err)
				client.Close()
				delete(scn.Clients, client)
			}
		}
	}
}
