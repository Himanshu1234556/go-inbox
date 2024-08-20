package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// Message represents a single chat message
type Message struct {
	UserID  string `json:"UserID"`
	Content string `json:"Content"`
	Type    string `json:"type,omitempty"` // Type is added to distinguish message types
}

// Global variables for managing state
var (
	chatHistory = []Message{}
	clients     = make(map[*websocket.Conn]string) // Map of websocket connections to usernames
	broadcast   = make(chan Message)
	upgrader    = websocket.Upgrader{}
	mu          sync.Mutex
)

func main() {
	// Setup CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5500", "http://127.0.0.1:5500"},
		AllowCredentials: true,
	})

	// Routes for WebSocket and chat history
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/chats", handleChatHistory)

	// Wrap the default mux with the CORS handler
	handler := corsHandler.Handler(http.DefaultServeMux)

	// Start handling messages
	go handleMessages()

	// Start the server
	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP connection to a WebSocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer ws.Close()

	// Register new client with their username
	var username struct {
		Username string `json:"username"`
	}
	err = ws.ReadJSON(&username)
	if err != nil {
		log.Println("Error while reading username:", err)
		return
	}

	mu.Lock()
	clients[ws] = username.Username
	mu.Unlock()

	// Notify all clients about the new user
	broadcastOnlineUsers()

	// Listen for incoming messages
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error while reading message:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			broadcastOnlineUsers()
			break
		}

		// Set the UserID to the username
		msg.UserID = username.Username

		// Store message in chat history
		mu.Lock()
		chatHistory = append(chatHistory, msg)
		mu.Unlock()

		// Broadcast the message to all clients
		msg.Type = "chat"
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error while sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func handleChatHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	json.NewEncoder(w).Encode(chatHistory)
	mu.Unlock()
}

func broadcastOnlineUsers() {
	var users []string
	mu.Lock()
	for _, username := range clients {
		users = append(users, username)
	}
	mu.Unlock()

	userListMessage := Message{
		Type:    "user_list",
		Content: "Online users",
	}

	mu.Lock()
	for client := range clients {
		err := client.WriteJSON(userListMessage)
		if err != nil {
			log.Println("Error while sending user list:", err)
			client.Close()
			delete(clients, client)
		}
	}
	mu.Unlock()
}
