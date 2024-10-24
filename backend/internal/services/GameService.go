package services

import (
	//"encoding/json"
	"dealer-backend/internal/models"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	//"net/http"
)

// StartMatchmaking handles the matchmaking logic for 4 players
func StartMatchmaking(players *models.PlayerConnections, playerID string) {
	// Main loop for matching players
	for {
		playerList := players.GetPlayerList()
		fmt.Println("connected players:", len(playerList))

		// Check if there are enough players to start a game
		if len(playerList) < 4 {
			time.Sleep(1 * time.Second) // Wait before checking again
			continue 
		}

		fmt.Println("enough players to start a game")

		// Create a map to hold selected players with their connections
		selectedPlayers := make(map[string]*websocket.Conn)
		conn, exists := players.GetPlayerConnection(playerID)
		if exists {
			selectedPlayers[playerID] = conn
		}

		// Select 3 random opponents
		for _, pid := range playerList {
			if pid != playerID {
				conn, exists := players.GetPlayerConnection(pid)
				if exists {
					selectedPlayers[pid] = conn
					if len(selectedPlayers) == 4 {
						break
					}
				}
			}
		}

		// Check if enough players were found
		if len(selectedPlayers) < 4 {
			time.Sleep(1 * time.Second) // Wait before checking again
			continue
		}

		// Create a slice of player IDs for the game struct
		playerIDs := make([]string, 0, 4)
		for pid := range selectedPlayers {
			playerIDs = append(playerIDs, pid)
		}

		// Create a new game
		gameID := fmt.Sprintf("game-%d", rand.Intn(100000)) // Generate a random game ID

		game := models.Game{
			GameID:  gameID,
			Players: playerIDs,
			State: models.GameState{
				Player1: models.Player{ID: playerIDs[0], Health: 100},
				Player2: models.Player{ID: playerIDs[1], Health: 100},
				Player3: models.Player{ID: playerIDs[2], Health: 100},
				Player4: models.Player{ID: playerIDs[3], Health: 100},
				Turn:    1, // Set initial turn to player 1
			},
		}

		fmt.Println("Starting chatService....")
		createRoom(&game, selectedPlayers)

		// Notify players about the new game
		fmt.Printf("Starting game %s with players: %v\n", gameID, playerIDs)

		break // Exit the loop after creating a game
	}
}
