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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	scd.Clients[conn] = true

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
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
				fmt.Println(err)
				client.Close()
				delete(scn.Clients, client)
			}
		}
	}
}
