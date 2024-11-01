package services

import (
	"dealer-backend/internal/models"
	"encoding/json"
	"fmt"
	"log"

	// "sync"
	"time"

	"github.com/gorilla/websocket"
)



func BroadcastAndAck(game *models.Game, connections map[string]*websocket.Conn, broadcastType string) {
	log.Println("Initializing acknowledgment listener...")

	log.Println("Broadcasting game state...")
	BroadcastGameState(game, connections, broadcastType)

	log.Println("Game state broadcasted, waiting for acknowledgments...")
	acknowledgmentReceived := WaitForAcks(ackChannel, connections)

	if acknowledgmentReceived {
		fmt.Println("All acknowledgments received for", broadcastType)
	} else {
		fmt.Println("No acknowledgment received for game state change. Timed out...")
	}
}


func WaitForAcks(ackChannel chan string, connections map[string]*websocket.Conn) bool {
    timeout := time.After(300 * time.Second)
    expectedAckCount := len(connections)
    acknowledgedPlayers := make(map[string]bool)
    ackCount := 0
    for {
        select {
        case playerID, ok := <-ackChannel:
            if !ok {
                fmt.Println("Ack channel closed")
                return ackCount == expectedAckCount
            }
            if !acknowledgedPlayers[playerID] {
                acknowledgedPlayers[playerID] = true
                ackCount++
                fmt.Println("Received acknowledgment from:", playerID)
            }
            if ackCount == expectedAckCount {
                return true
            }
        case <-timeout:
            var nonAcknowledgedPlayers []string
            for playerID := range connections {
                if !acknowledgedPlayers[playerID] {
                    nonAcknowledgedPlayers = append(nonAcknowledgedPlayers, playerID)
                }
            }
            fmt.Println("Acknowledgment timeout. Players who did not acknowledge:", nonAcknowledgedPlayers)
            return false
        }
    }
}



func BroadcastGameState(game *models.Game, connections map[string]*websocket.Conn, stateType string) {
	
	var message map[string]interface{}
    var currentPlayer models.Player

    // Determine the current player based on the turn
    switch game.State.Turn {
        case 1:
            currentPlayer = game.State.Player1
        case 2:
            currentPlayer = game.State.Player2
        case 3:
            currentPlayer = game.State.Player3
        case 4:
            currentPlayer = game.State.Player4
        default:
            // Handle cases where Turn doesn't match (optional)
            currentPlayer = models.Player{}
        }
	// Assuming game is of type *models.Game
	if stateType == "gamestate" {
		message = map[string]interface{}{
			"type": stateType,
			"data": game, // Send the full game object if stateType is "gamestate"
		}
	} else if stateType == "healthstate" {
		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"player": currentPlayer.ID,  // Use current player ID
				"health": currentPlayer.Health,  // Use current player Health
			},
		}
	}else if stateType == "cardplayed" {
		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"playerId": currentPlayer.ID,  // Use current player ID
				"card": currentPlayer.PlayedCard,  // Use current player Health
			},
		}
    } else if stateType == "trickwon" {
		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"player": game.State.RoundWinner,  // Use current player ID
				"score": game.State.RoundWinner.Score,
			}, 
		}
    }else if stateType == "resetcardplayed" {
		message = map[string]interface{}{
			"type": stateType,
		}
    }else if stateType == "biddingcomplete" {
		message = map[string]interface{}{
			"type": stateType,
		}
    }else if stateType == "bidupdate" {
		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"playerId": currentPlayer.ID,
				"bid": currentPlayer.Bid,
			}, 	
		}
    }else if stateType == "gameover" {
		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"playerId": currentPlayer.ID,

			}, 	
		}
    }
	
	
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		// Handle error (e.g., log it)
		return
	}

	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			// Handle error (e.g., log it or remove the connection)
		}
	}
}


