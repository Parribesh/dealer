package services

import (
	//"encoding/json"
	"dealer-backend/internal/models"
	"fmt"

	// "log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message represents a chat message structure
type Message struct {
    SenderID string `json:"sender_id"`
    Content  string `json:"content"`
}

// StartChatService handles chat communication between two players
func StartChatService(game models.Game, playerConnections *models.PlayerConnections) {
    // Retrieve connections for both players
    player1Conn, exists1 := playerConnections.GetPlayerConnection(game.Players[0])
    player2Conn, exists2 := playerConnections.GetPlayerConnection(game.Players[1])

    fmt.Printf("Both connection recieved.. Starting Reading GoRoutines")
    
    if !exists1 || !exists2 {
        fmt.Println("One of the player connections is missing. Exiting chat service.")
        return
    }

	

    // Use a WaitGroup to manage goroutines for reading messages
    var wg sync.WaitGroup

    // Function to handle message reading and sending
    readMessages := func(conn *websocket.Conn, otherConn *websocket.Conn, playerID string) {
        defer wg.Done() // Notify when done

		fmt.Println("Listening for player", playerID)
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                fmt.Println("Error reading message:", err)
                break // Exit if there's an error
            }

            // Construct the message object to send
            message := Message{
                SenderID: playerID,
                Content:  string(msg),
            }


            // Send the message to the other player
            if err := otherConn.WriteJSON(message); err != nil {
                fmt.Println("Error sending message:", err)
                break // Exit if there's an error
            }
        }
    }

    wg.Add(2) // We have two players

    // Start reading messages from player 1 and send to player 2
    go readMessages(player1Conn, player2Conn, game.Players[0])
    // Start reading messages from player 2 and send to player 1
    go readMessages(player2Conn, player1Conn, game.Players[1])

    // Wait for both goroutines to finish
    wg.Wait()

    // Close connections after chat is done
    player1Conn.Close()
    player2Conn.Close()
}
