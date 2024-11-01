package services

import (
	"dealer-backend/internal/models"
	// "encoding/json"
	"fmt"
	"log"
	// "sync"
	"time"

	"github.com/gorilla/websocket"
)


// type BidMessage struct {
// 	Type string `json:"type"`
// 	PlayerID string `json:"playerId"`
// 	Bid      int    `json:"bid"`

// }


func WaitForAllBids(game *models.Game, connections map[string]*websocket.Conn) map[string]int {
	
	timeout := time.After(120 * time.Second)
	expectedBidCount := len(connections)

	bids := make(map[string]int)
	acknowledgedPlayers := make(map[string]bool)

	fmt.Println("Starting to listen for bids from all players...")

	// Start the main loop to collect bids or timeout in a separate goroutine
	
	func() {
		bidCount := 0
		ticker := time.NewTicker(60 * time.Second) // Create a ticker that ticks every 10 seconds
		defer ticker.Stop() // Ensure the ticker is stopped when done

		for {
			select {
			case bid, ok := <-bidChannel:
				if !ok {
					fmt.Println("Bid channel closed.")
					return
				}

				SetPlayerBid(game, bid.PlayerID, bid.Bid)
				bids[bid.PlayerID] = bid.Bid
				acknowledgedPlayers[bid.PlayerID] = true
				bidCount++
				fmt.Printf("Processed bid from player %s: %d\n", bid.PlayerID, bid.Bid)

				// Send bid update notification
				BroadcastAndAck(game, connections, "gamestate")

				if bidCount == expectedBidCount {
					fmt.Println("All bids received. Broadcasting bidding complete.")
					BroadcastGameState(game, connections, "biddingcomplete")
					BroadcastAndAck(game, connections, "gamestate")
					return
				}
			case <-timeout:
				fmt.Println("Bidding timeout. Players who did not bid:")
				for playerID := range connections {
					if !acknowledgedPlayers[playerID] {
						fmt.Printf("Missing bid from player: %s\n", playerID)
					}
				}
				return
			case <-ticker.C: // Listen for ticker ticks
				// Send a reminder message only to clients who have not submitted bids
				for playerID, conn := range connections {
					if !acknowledgedPlayers[playerID] { // Check if the player has not submitted a bid
						msg := map[string]interface{}{
							"type": "updatebid",
							"data": map[string]string{
								"message": "Please update your bid.",
							},
						}
						if err := conn.WriteJSON(msg); err != nil {
							log.Printf("Error sending update bid message to player %s: %v\n", playerID, err)
						} else {
							fmt.Printf("Sent update bid message to player %s.\n", playerID)
						}
					}
				}
			}
		}
	}()
				
	
	return bids
}

func SetPlayerBid(game *models.Game, playerID string, bidAmount int) error {
    // Ensure the Bids map is initialized
    if game.State.Bids == nil {
        game.State.Bids = make(map[string]int)
    }

    // Set the bid in the Bids map for the player
    game.State.Bids[playerID] = bidAmount

    // Update the Player struct's Bid field
    switch playerID {
    case game.State.Player1.ID:
        game.State.Player1.Bid = bidAmount
    case game.State.Player2.ID:
        game.State.Player2.Bid = bidAmount
    case game.State.Player3.ID:
        game.State.Player3.Bid = bidAmount
    case game.State.Player4.ID:
        game.State.Player4.Bid = bidAmount
    default:
        return fmt.Errorf("player with ID %s not found", playerID)
    }

    return nil
}



