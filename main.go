package main

import (
	"log"
	"net/http"

	"andreswang.com/dialoq/internal"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo purposes
	},
}

// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade the HTTP connection to a WebSocket connection
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Error upgrading to WebSocket: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Simple echo functionality
// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading message: %v", err)
// 			return
// 		}
// 		// Echo the message back to the client
// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Printf("Error writing message: %v", err)
// 			return
// 		}
// 	}
// }

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", internal.HandleWebSocket)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	log.Printf("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
