package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dealer-backend/internal/auth"
	"dealer-backend/internal/config"
	"dealer-backend/internal/handlers"
	"dealer-backend/internal/middlewares"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all connections
    },
}


// WebSocket handler for player connections
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	fmt.Println("WebSocket connection established")

	// Declare tokenMsg outside of the loop for use later
	var tokenMsg []byte

	// Message reading loop (wait for messages from the client)
	for {
		// Read the token message from the client
		_, tokenMsg, err = conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading token message:", err)
			conn.Close()
			return
		}

        if(len(tokenMsg) > 0){
            // Process the token message
            fmt.Println("Token received:", string(tokenMsg))
            break
        }	  
	}

	// Parse the token message
	var tokenData struct {
		Token string `json:"token"`
	}

	// Use the tokenMsg (received in the loop) for parsing
	if err := json.Unmarshal(tokenMsg, &tokenData); err != nil {
		fmt.Println("Error parsing token message:", err)
		conn.Close()
		return
	}

	// Extract username (playerID) from the token
	username, err := auth.GetUsernameFromToken(tokenData.Token)
	fmt.Println("Username from token: ", username)
	if err != nil {
		fmt.Println("Error extracting username from token:", err)
		conn.Close()
		return
	}

	// Add the player connection to the map
	config.PlayerConnections.AddPlayer(username, conn)

	// Broadcast that the user has joined
	config.PlayerConnections.BroadcastMessage("message", username+" joined")

	// Optionally start matchmaking or other services
	// services.StartMatchmaking(playerConnections, username)
}


// Handler function for the root route
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Card Game API")
}

// Main function to start the web server
func main() {
    // Define routes
    http.HandleFunc("/", homePage)
    http.HandleFunc("/protected", handlers.ProtectedHandler)

    //User Authentication endpoings
    http.HandleFunc("/login", handlers.LoginHandler)
    //http.HandleFunc("/register", handlers.RegisterHandler)

    //Lobby-Matchmaking handler
    http.HandleFunc("/lobby/join", handlers.JoinHandler)
    http.HandleFunc("/lobby/status", handlers.LobbyStatusHandler)

    //Websocker handler 
    http.HandleFunc("/ws", wsHandler) // WebSocket route

    //Game handler
    http.Handle("/game/start", auth.ValidateJWTMiddleware(http.HandlerFunc(handlers.StartHandler)))

    http.HandleFunc("/game/move", handlers.MoveHandler)
    //http.HandleFunc("/game/draw", handlers.DrawHandler)
    //http.HandleFunc("/game/status/", handlers.GameStatusHandler)
    //http.HandleFunc("/game/result/",handlers.ResultHandler)
    //http.HandleFunc("/game/history/", handlers.ResultHandler)


    // Start the server on port 8080
    log.Println("Starting server on port 8080...")
    // Start the server and handle any errors
    if err := http.ListenAndServe(":8080", middlewares.CorsMiddleware(http.DefaultServeMux)); err != nil {
        log.Fatalf("Server failed to start: %v", err) // Log the error and exit the program
    }

}
