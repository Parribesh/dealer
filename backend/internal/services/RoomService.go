package services

import (
	"dealer-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func createRoom(game *models.Game, connections map[string]*websocket.Conn) {
	game.ShuffleAndDealCards()

	// Start the game loop in a separate goroutine
	go gameLoop(game, connections)

	// Send initial game state
	broadcastGameState(game, connections, "gamestate")
}

func gameLoop(game *models.Game, connections map[string]*websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Println("Game loop started")

	for range ticker.C {
		fmt.Println("Game loop ticked")
		currentPlayerNumber := game.State.Turn
		fmt.Println("Current player number:", currentPlayerNumber)
		currentPlayer := getCurrentPlayer(&game.State, currentPlayerNumber)
		fmt.Println("Current player:", currentPlayer)

		//currentPlayer.Health -= currentPlayer.Health/60
		currentPlayer.Health -= 10
		fmt.Println("Current player health:", currentPlayer.Health)
		fmt.Println("Broadcasting game state")

		broadcastGameState(game, connections, "healthstate")
		fmt.Println("Game state broadcasted")
		if currentPlayer.Health <= 0 {
			// Move to the next player with health > 0
				game.State.Turn = (game.State.Turn + 1) % 5 
				if( game.State.Turn == 0) {
					fmt.Println("Broadcasting game state when health is 0")
					broadcastGameState(game, connections, "healthstate")
					break
				}
				nextPlayer := getCurrentPlayer(&game.State, (game.State.Turn ))
				if nextPlayer.Health <= 0  {
					fmt.Println("currnt turn: ", game.State.Turn)
					break
				}
		}


		// Check if the game is over (only one player with health > 0)
		// if isGameOver(&game.State) {
		// 	// Handle game over logic here
		// 	break
		// }
	}
}

func broadcastGameState(game *models.Game, connections map[string]*websocket.Conn, stateType string) {
	
	var message map[string]interface{}

	// Assuming game is of type *models.Game
	if stateType == "gamestate" {
		message = map[string]interface{}{
			"type": stateType,
			"data": game, // Send the full game object if stateType is "gamestate"
		}
	} else if stateType == "healthstate" {
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

		message = map[string]interface{}{
			"type": stateType,
			"data": map[string]interface{}{
				"player": currentPlayer.ID,  // Use current player ID
				"health": currentPlayer.Health,  // Use current player Health
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

func getCurrentPlayer(state *models.GameState, playerNumber int) *models.Player {
	switch playerNumber {
	case 1:
		return &state.Player1
	case 2:
		return &state.Player2
	case 3:
		return &state.Player3
	case 4:
		return &state.Player4
	default:
		return nil
	}
}

func isGameOver(state *models.GameState) bool {
	playersAlive := 0
	if state.Player1.Health > 0 {
		playersAlive++
	}
	if state.Player2.Health > 0 {
		playersAlive++
	}
	if state.Player3.Health > 0 {
		playersAlive++
	}
	if state.Player4.Health > 0 {
		playersAlive++
	}
	return playersAlive <= 1
}
