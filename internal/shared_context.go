package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	// "time"

	"github.com/gorilla/websocket"
)

var (
	clients     = make(map[*websocket.Conn]bool)
	clientMutex sync.RWMutex

	globalSession      sessionMessages
	globalSessionMutex = &sync.RWMutex{}
)

type singleMessage struct {
	M        string `json:"m"`        // message text
	Username string `json:"username"` // sender name
	// Timestamp time.Time `json:"timestamp"` // optionally add a timestamp
}

type sessionMessages struct {
	Messages []singleMessage `json:"messages"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: check the origin
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	clientMutex.Lock()
	clients[conn] = true
	clientMutex.Unlock()

	defer func() {
		clientMutex.Lock()
		delete(clients, conn)
		clientMutex.Unlock()
		conn.Close()
	}()

	err = sendAllSessionMessages(conn)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		if messageType == websocket.TextMessage {
			log.Println("Received message:", string(message))

			singleMessage := singleMessage{
				M:        string(message),
				Username: "default",
			}
			addToCurrentSession(singleMessage)

			messageJson, err := json.Marshal(globalSession)
			if err != nil {
				log.Printf("Error marshaling existing messages: %v", err)
				return
			}
			err = broadcastMessage(messageJson)
			if err != nil {
				log.Printf("Error marshaling existing messages: %v", err)
				return
			}
		}
	}
}

func broadcastMessage(messageJson []byte) error {
	for conn, v := range clients {
		if v {
			err := conn.WriteMessage(websocket.TextMessage, messageJson)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return err
			}
		}
	}
	return nil
}

func addToCurrentSession(text singleMessage) {
	globalSessionMutex.Lock()
	defer globalSessionMutex.Unlock()

	globalSession.Messages = append(globalSession.Messages, text)
}

func sendAllSessionMessages(conn *websocket.Conn) error {
	data, err := json.Marshal(globalSession)
	if err != nil {
		log.Printf("Error marshaling existing messages: %v", err)
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}
	return nil
}
