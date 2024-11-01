package services

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type BidMessage struct {
	Type     string `json:"type"`
	PlayerID string `json:"playerId"`
	Bid      int    `json:"bid,omitempty"` // Optional field for bid
}

var (
	bidChannel = make(chan BidMessage, 8)
	ackChannel  = make(chan string, 8)
)


func StartMessageRouter(connections map[string]*websocket.Conn) {
	
	// Start a goroutine to handle each player's connection and route messages
	for playerID, conn := range connections {
		go func(playerID string, conn *websocket.Conn) {
			defer conn.Close()
			log.Println("Waiting for messagess..")
			for {
				// Read message from the WebSocket connection
				_, rawMessage, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading message from player %s: %v\n", playerID, err)
					return
				}


				log.Println("rawMessage received..", rawMessage)

				// Unmarshal message into a BidMessage struct
				var msg BidMessage
				if err := json.Unmarshal(rawMessage, &msg); err != nil {
					log.Printf("Failed to unmarshal message from player %s: %v\n", playerID, err)
					continue
				}

				// Route messages based on the Type field
				switch msg.Type {
				case "placebid":
					// Set the player ID and send to bidChannel
					msg.PlayerID = playerID
					bidChannel <- msg

				case "acknowledgment":
					// Send player ID to ackChannel
					ackChannel <- playerID

				default:
					log.Printf("Unknown message type from player %s: %v\n", playerID, msg.Type)
				}
			}
		}(playerID, conn)
	}

}
