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

var gameRooms = make(map[string]*models.Game)

func createRoom(game *models.Game, connections map[string]*websocket.Conn) {
	gameRooms[game.GameID] = game  // Store game instance by its ID
	game.ShuffleAndDealCards()
	StartMessageRouter(connections)
	// Send initial game state
	BroadcastAndAck(game, connections, "gamestate")
	// Start the game loop in a separate goroutine
	go gameLoop(game, connections)

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


// ************************** CARD LOGIC ********************************************
// compareCards compares two cards, giving precedence to suits and then ranks
func compareCards(card1, card2 *models.Card) int {
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
// ************************** CARD LOGIC - END ********************************************




// ************************** MOVE LOGIC ********************************************
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

    // Find the card in the player's hand to set as PlayedCard
    for i, card := range currentPlayer.Hand {
        if card.Suit == playedCardMsg.Card.Suit && card.Rank == playedCardMsg.Card.Rank {
            // Assign the address of the existing card in hand to PlayedCard
            currentPlayer.PlayedCard = &currentPlayer.Hand[i]
            fmt.Println("Player", playerID, "played a card:", currentPlayer.PlayedCard)
            return
        }
    }
}

// resetPlayedCards resets the PlayedCard for all players in the game state
func resetPlayedCards(game *models.Game, connections map[string]*websocket.Conn) {
	game.State.Player1.RemovePlayedCard()
    game.State.Player2.RemovePlayedCard()
    game.State.Player3.RemovePlayedCard()
    game.State.Player4.RemovePlayedCard()
	game.State.TrickSuit = ""
	BroadcastAndAck(game, connections, "resetcardplayed" )
}


// **********************************MOVE LOGIC - END *********************************


// ********************************** ROOM STATE LOGIC - END *********************************

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


func determineTrickWinner(state *models.GameState) *models.Player {
    var highestTrumpCard *models.Card
    var highestSpadeCard *models.Card
    // var highestOtherCard *models.Card
    var winner *models.Player
    var trumpSuit = state.TrickSuit

    for _, player := range []*models.Player{&state.Player1, &state.Player2, &state.Player3, &state.Player4} {
        if player.PlayedCard != nil {
            playedCard := player.PlayedCard

            if playedCard.Suit == trumpSuit && highestSpadeCard == nil {
                // Track the highest card in the trump suit
                if highestTrumpCard == nil || compareCards(playedCard, highestTrumpCard) > 0 {
                    highestTrumpCard = playedCard
                    winner = player
                }
            } else if playedCard.Suit == "S" {
                log.Println("Turup played... ")
                // Track the highest spade card if no trump suit is played
                if highestSpadeCard == nil || compareCards(playedCard, highestSpadeCard) > 0 {
                log.Println("Highest spade card set, winner till now.. ", player.ID)
                    highestSpadeCard = playedCard
                    winner = player
                }
            } 
            // else {
            //     // Track the highest card in other suits if no trump suit or spades are played
            //     if highestOtherCard == nil || compareCards(playedCard, highestOtherCard) > 0 {
            //         highestOtherCard = playedCard
            //         winner = player
            //     }
            // }
        }
    }

    // Final decision: priority order spades > trump > other
    if highestSpadeCard != nil {
        return winner // highest spade card is the winner
    } else if highestTrumpCard != nil {
        return winner //  trump card is the winner if no trump suit is played
    } else {
        return winner // otherwise, the highest other card is the winner
    }
}


// updateGameStateAfterTrick updates the game state after a trick is completed
func updateGameStateAfterTrick(game *models.Game, winner *models.Player, connections map[string]*websocket.Conn) {
    // Update the round winner in the GameState
    game.State.RoundWinner = winner
	winner.Score +=1
    // Increment the winner's score
    if game.State.Scores == nil {
        game.State.Scores = make(map[string]int)
    }
    
    game.State.Scores[winner.ID]++

    // Optional: Reset played cards for the next trick
    resetPlayedCards(game, connections)

    // Debug log to confirm the updates
    fmt.Printf("Updated game state: round winner is %s, new score is %d\n", winner.ID, game.State.Scores[winner.ID])
}

func isGameOver(state *models.GameState) bool {
    // Check each player's hand for emptiness
    for _, player := range []*models.Player{&state.Player1, &state.Player2, &state.Player3, &state.Player4} {
        if len(player.Hand) == 0 { // Assuming Hand is a slice of cards
            return true // Game is over if any player's hand is empty
        }
    }
    return false // Game continues if all players have cards in their hands
}





//*************************** MAIN LOOP ***********************************

func gameLoop(game *models.Game, connections map[string]*websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Println("Game loop started")

	// Players := []*models.Player{
	// 	&game.State.Player1,
	// 	&game.State.Player2,
	// 	&game.State.Player3,
	// 	&game.State.Player4,
		
	// }
    
	WaitForAllBids(game, connections);
	
	// Main Game Loop
	for range ticker.C {
		fmt.Println("Game loop ticked")
		currentPlayerNumber := game.State.Turn
		currentPlayer := getCurrentPlayer(&game.State, currentPlayerNumber)

		// Check if the player's health has dropped (indicating timer expiration)
		currentPlayer.Health -= 1
		fmt.Println("Current player health:", currentPlayer.Health)
		BroadcastGameState(game, connections, "healthstate")
		if currentPlayer.Health <= 0 {
			// Timeout: move to the next player if the current player did not play a card
			advanceTurn(game)
			resetHealthForNextTurn(game)
			continue
		}

		// Wait for the current player to play a card
		if currentPlayer.PlayedCard == nil {
			continue
		}

		// Process the played card
		processPlayedCard(game, currentPlayer)

		// Validate the played card
		if isValidCard(currentPlayer, game.State.TrickSuit, game.State) {
			BroadcastGameState(game, connections, "cardplayed")
		} else {
			fmt.Println("Invalid card played")
			continue
		}

        hasWinner := false
		// Check if all players have played (end of trick)
		if allPlayersHavePlayed(game) {
			winner := determineTrickWinner(&game.State)
            hasWinner = true
			updateGameStateAfterTrick(game, winner, connections)
			BroadcastAndAck(game, connections, "trickwon")
			BroadcastGameState(game, connections, "gamestate")
            game.State.Turn = game.State.GetPlayerPosition(*winner)
		}

		// Move to the next turn
        if !hasWinner {
		    advanceTurn(game)
        }
        if isGameOver(&game.State) {
            BroadcastGameState(game, connections, "gameover")
            log.Println("Game Over!! Thank you for playing...")
            return
        }
		resetHealthForNextTurn(game)
	}
}



