package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"sd-chat/domain/services"
	"sd-chat/infrastructure/config"
	"sd-chat/web/controllers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults instead.\n")
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Error loading configs...\n")
	}
	// Initializing Server Connection Data (for handling clients and messages)
	scn := web_socket_service.NewServerConnData()
	http.HandleFunc("/", web.ServeHome)

	// Serve all CSS and JS
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir("frontend/css"))))

	http.Handle("/js/", http.StripPrefix("/js/",
		http.FileServer(http.Dir("frontend/js"))))

	http.HandleFunc("/ws", scn.HandleConnections)

	go scn.HandleMessages()

	fmt.Printf("Server started on port %s\n", cfg.ServerPort)
	err = http.ListenAndServe(cfg.ServerPort, nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
