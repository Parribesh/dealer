package services

import (
	"dealer-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var gameRooms = make(map[string]*models.Game)

func createRoom(game *models.Game, connections map[string]*websocket.Conn) {
	gameRooms[game.GameID] = game  // Store game instance by its ID
	game.ShuffleAndDealCards()
	
	// Start the game loop in a separate goroutine
	go gameLoop(game, connections)

	// Send initial game state
	broadcastGameState(game, connections, "gamestate")
}

func resetHealthForNextTurn(game *models.Game) {
    currentTurn := game.State.Turn

    // If the turn is 0, set it to 1
    if currentTurn == 0 {
        currentTurn = 1
    }

    nextPlayer := getCurrentPlayer(&game.State, currentTurn)

    // Reset the health of the next player (i.e., the player's timer)
    nextPlayer.Health = 100 // or set this to the maximum health/timer value

    fmt.Println("Health reset for player", currentTurn, "to", nextPlayer.Health)
}

func advanceTurn(game *models.Game) {
    // Advance the turn, making sure 0 becomes 1
    game.State.Turn = (game.State.Turn + 1) % 5

    if game.State.Turn == 0 {
        game.State.Turn = 1 // If the turn reaches 0, set it to 1
    }

    resetHealthForNextTurn(game)
}

func isValidCard(player *models.Player, trickSuit models.Suit, gameState models.GameState) bool {
    // If the player has a card in hand with the TrickSuit, they must play it
    hasTrickSuit := false
    for _, card := range player.Hand {
        if card.Suit == trickSuit {
            hasTrickSuit = true
            break
        }
    }

    // If the player has the trick suit, the card they play must match the trick suit
    if hasTrickSuit && player.PlayedCard.Suit != trickSuit {
        fmt.Println("Player must follow suit")
        return false
    }

    // If the player doesn't have the trick suit, they can play any card
    return true
}

func processPlayedCard(game *models.Game, currentPlayer *models.Player) {
    // If it's the first card of the trick, set the trick suit
    if game.State.TrickSuit == "" {
        game.State.TrickSuit = currentPlayer.PlayedCard.Suit
        fmt.Println("Trick suit set to", game.State.TrickSuit)
    }

    // Apply other game logic for the played card
}

func allPlayersHavePlayed(game *models.Game) bool {
    // Check if all players have played a card
    if game.State.Player1.PlayedCard == nil {
        return false
    }
    if game.State.Player2.PlayedCard == nil {
        return false
    }
    if game.State.Player3.PlayedCard == nil {
        return false
    }
    if game.State.Player4.PlayedCard == nil {
        return false
    }

    return true
}



func HandlePlayerMove(gameID string, playerID string, message []byte) {
	game := gameRooms[gameID]
    var playedCardMsg models.PlayedCardMessage
    
    // Unmarshal the incoming message to get the played card details
    err := json.Unmarshal(message, &playedCardMsg)
    if err != nil {
        fmt.Println("Error unmarshalling played card message:", err)
        return
    }

    // Find the current player by their ID
    currentPlayer := getPlayerByID(game, playerID)
    if currentPlayer == nil {
        fmt.Println("Player not found:", playerID)
        return
    }

    // Update the current player's PlayedCard
    currentPlayer.PlayedCard = &models.Card{
        Suit: playedCardMsg.Card.Suit,
        Rank: playedCardMsg.Card.Rank,
    }

    fmt.Println("Player", playerID, "played a card:", currentPlayer.PlayedCard)
}

func getPlayerByID(game *models.Game, playerID string) *models.Player {
    if game.State.Player1.ID == playerID {
        return &game.State.Player1
    }
    if game.State.Player2.ID == playerID {
        return &game.State.Player2
    }
    if game.State.Player3.ID == playerID {
        return &game.State.Player3
    }
    if game.State.Player4.ID == playerID {
        return &game.State.Player4
    }
    return nil
}


// compareCards compares two cards, giving precedence to suits and then ranks
func compareCards(card1, card2 *models.Card) int {
    // Compare suits first
    if models.SuitRankings[card1.Suit] > models.SuitRankings[card2.Suit] {
        return 1 // card1 is higher
    } else if models.SuitRankings[card1.Suit] < models.SuitRankings[card2.Suit] {
        return -1 // card2 is higher
    }

    // If suits are equal, compare ranks
    return compareRanks(card1.Rank, card2.Rank)
}

// compareRanks compares ranks of two cards
func compareRanks(rank1, rank2 models.Rank) int {
    rankOrder := map[models.Rank]int{
        "2":  0,
        "3":  1,
        "4":  2,
        "5":  3,
        "6":  4,
        "7":  5,
        "8":  6,
        "9":  7,
        "10": 8,
        "J":  9,
        "Q":  10,
        "K":  11,
        "A":  12,
    }

    return rankOrder[rank1] - rankOrder[rank2]
}

func determineTrickWinner(state *models.GameState) *models.Player {
    var highestCard *models.Card
    var winner *models.Player
    var trumpSuit = state.TrickSuit

    trickSuitPlayed := false

    for _, player := range []*models.Player{&state.Player1, &state.Player2, &state.Player3, &state.Player4} {
        if player.PlayedCard != nil {
            playedCard := player.PlayedCard

            // Check if the played card is of the trick suit
            if playedCard.Suit == trumpSuit {
                trickSuitPlayed = true

                // Compare using suit rankings and then rank
                if highestCard == nil || compareCards(playedCard, highestCard) > 0 {
                    highestCard = playedCard
                    winner = player
                }
            } else {
                // Handle fallback suit logic if trick suit is not played
                if !trickSuitPlayed {
                    if highestCard == nil || compareCards(playedCard, highestCard) > 0 {
                        highestCard = playedCard
                        winner = player
                    }
                }
            }
        }
    }

    return winner
}

// updateGameStateAfterTrick updates the game state after a trick is completed
func updateGameStateAfterTrick(game *models.Game, winner *models.Player) {
    // Update the round winner in the GameState
    game.State.RoundWinner = winner

    // Increment the winner's score
    if game.State.Scores == nil {
        game.State.Scores = make(map[string]int)
    }
    
    game.State.Scores[winner.ID]++

    // Optional: Reset played cards for the next trick
    resetPlayedCards(&game.State)

    // Debug log to confirm the updates
    fmt.Printf("Updated game state: round winner is %s, new score is %d\n", winner.ID, game.State.Scores[winner.ID])
}

// resetPlayedCards resets the PlayedCard for all players in the game state
func resetPlayedCards(state *models.GameState) {
    state.Player1.PlayedCard = nil
    state.Player2.PlayedCard = nil
    state.Player3.PlayedCard = nil
    state.Player4.PlayedCard = nil
}


func waitForAllAcknowledgments(connections map[string]*websocket.Conn) bool {
	ackChannel := make(chan string) // Channel to collect acknowledgments
	timeout := time.After(5 * time.Second) // Set a timeout for acknowledgment
	expectedAckCount := len(connections)

	// Create a slice to keep track of all player IDs
	playerIDs := make([]string, 0, expectedAckCount)
	for playerID := range connections {
		playerIDs = append(playerIDs, playerID) // Collecting expected player IDs
	}

	acknowledgedPlayers := make(map[string]bool) // Map to track acknowledged players

	go func() {
		for _, conn := range connections {
			// Use an inline struct definition for PlayerMessage
			var msg struct {
				Type     string `json:"type"`
				PlayerID string `json:"playerId"`
			}

			if err := conn.ReadJSON(&msg); err == nil && msg.Type == "acknowledgment" {
				ackChannel <- msg.PlayerID // Send the PlayerID of the acknowledging player
			}
		}
	}()

	// Track the number of acknowledgments received
	ackCount := 0
	for {
		select {
		case playerID := <-ackChannel:
			ackCount++
			acknowledgedPlayers[playerID] = true // Mark this player as acknowledged
			fmt.Println("Received acknowledgment from:", playerID)
			// Check if we have received acknowledgments from all players
			if ackCount == expectedAckCount {
				return true
			}
		case <-timeout:
			// Check which players did not acknowledge
			var nonAcknowledgedPlayers []string
			for _, playerID := range playerIDs {
				if !acknowledgedPlayers[playerID] {
					nonAcknowledgedPlayers = append(nonAcknowledgedPlayers, playerID)
				}
			}
			fmt.Println("Acknowledgment timeout. Players who did not acknowledge:", nonAcknowledgedPlayers)
			return false // Return false if timeout occurs
		}
	}
}



func gameLoop(game *models.Game, connections map[string]*websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Println("Game loop started")

	
	for range ticker.C {
		fmt.Println("Game loop ticked")
		currentPlayerNumber := game.State.Turn
		currentPlayer := getCurrentPlayer(&game.State, currentPlayerNumber)
		
		// Check if the player's health has dropped (indicating timer expiration)
		currentPlayer.Health -= 1
		fmt.Println("Current player health:", currentPlayer.Health)
		fmt.Println("Broadcasting game state")
		broadcastGameState(game, connections, "healthstate")
		fmt.Println("Game state broadcasted")
		if currentPlayer.Health <= 0 {
			// Timeout: move to the next player if the current player did not play a card
			advanceTurn(game)
			resetHealthForNextTurn(game) // Reset health for the next player
			continue
		}

		// Allow the current player to play a card (wait for user input or API call)
		if currentPlayer.PlayedCard == nil {
			// Wait for the player to play a card via frontend interaction or websocket
			continue // Skip to next loop iteration until the player plays a card
		}

		// Process the card played by the current player
        processPlayedCard(game, currentPlayer)

		// Validate the played card against the rules
		if isValidCard(currentPlayer, game.State.TrickSuit, game.State) {

            
			// Broadcast the updated game state with the played card
			broadcastGameState(game, connections, "gamestate")
            // Wait for acknowledgment for the game state change
			acknowledgmentReceived := waitForAllAcknowledgments(connections)
			if acknowledgmentReceived {
				fmt.Println("Acknowledgment received for game state change.")
				
				// Now broadcast the card played message
				broadcastGameState(game, connections, "cardplayed")
			} else {
				fmt.Println("No acknowledgment received for game state change. Handling timeout...")
				// Handle the case where acknowledgment is not received
				// You might want to prompt the player again or log the incident
				continue // Optionally skip the cardplayed broadcast if no ack received
			}


		} else {
			// Handle invalid card case (e.g., reset card or prompt the player again)
			fmt.Println("Invalid card played")
			continue
		}


		// Check if all players have played a card (end of trick)
		if allPlayersHavePlayed(game) {
			// Determine the winner of the trick based on the leading suit and trump cards
			winner := determineTrickWinner(&game.State)
			// Update game state, increment the winner's trick count
			updateGameStateAfterTrick(game, winner)
			broadcastGameState(game, connections, "gamestate")
            acknowledgmentReceived := waitForAllAcknowledgments(connections)
			if acknowledgmentReceived {
				fmt.Println("Acknowledgment received for game state change.")
				
				// Now broadcast the card played message
				broadcastGameState(game, connections, "trickwon")
			} else {
				fmt.Println("No acknowledgment received for game state change. Handling timeout...")
				// Handle the case where acknowledgment is not received
				// You might want to prompt the player again or log the incident
				continue // Optionally skip the cardplayed broadcast if no ack received
			}
			
		}

		// Move to the next player's turn
		advanceTurn(game)
		resetHealthForNextTurn(game)
	}
}


func broadcastGameState(game *models.Game, connections map[string]*websocket.Conn, stateType string) {
	
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
