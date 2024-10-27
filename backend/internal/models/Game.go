package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"encoding/json"
	"math/rand"

	"github.com/gorilla/websocket"
)

// Game stores the state of an active game
type Game struct {
	GameID    string
	Players   []string // List of player IDs
	State     GameState // Game state can be complex or simplified
	GameWinner *Player
}

type GameState struct {
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Player3 Player `json:"player3"`
	Player4 Player `json:"player4"`
	Turn    int    `json:"turn"` // Keeps track of which player's turn it is
	TrickSuit Suit `json:"trick_suit"`
	RoundWinner *Player `json:"round_winner,omitempty"`
	Scores      map[string]int   `json:"scores"` // Track scores by player ID
	Bids        map[string]int   `json:"bids"` 
}

type Player struct {
	ID    string  `json:"id"`
	Hand  []Card  `json:"hand"`
	Health int    `json:"health"`
	PlayedCard *Card  `json:"played_card"`
	Bid      int    `json:"bid"` // Number of tricks the player aims to win
}

// RemovePlayedCard removes the played card from the player's hand
func (p *Player) RemovePlayedCard() {
	fmt.Println("Remove requested....")
    if p.PlayedCard == nil {
        fmt.Println("No played card to remove.")
        return
    }

   fmt.Printf("Played card before removal: %+v, Address: %p\n", p.PlayedCard, p.PlayedCard)

    // Find the index of PlayedCard in Hand
    for i := range p.Hand {
        // Compare the values of the card
        if p.Hand[i].Suit == p.PlayedCard.Suit && p.Hand[i].Rank == p.PlayedCard.Rank {
            // Capture the card to log it before removal
            removedCard := p.Hand[i]

            // Remove the card at index i
            p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
            fmt.Printf("Removed %s from hand.\n", removedCard)

            // Optionally reset PlayedCard to nil
            p.PlayedCard = nil // Reset the pointer to nil
            return
        }
    }

    fmt.Println("Played card not found in hand.")
}



type Suit string
type Rank string

const (
	Hearts   Suit = "H"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
	Spades   Suit = "S"
)

const (
	Ace   Rank = "A"
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

// Define a map to quickly get the suit ranking
var SuitRankings = map[Suit]int{
    Hearts:   0,
    Diamonds: 1,
    Clubs:    2,
    Spades:   3,
}

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

// Method to get the card's identifier (e.g., "6D" for Six of Diamonds)
func (c Card) Identifier() string {
	return string(c.Rank) + string(c.Suit)
}

type PlayedCardMessage struct {
    PlayerID string `json:"player_id"`
    Card     Card `json:"card"`
 
}


// Struct to hold the player connections
type PlayerConnections struct {
	mu      sync.RWMutex
	players map[string]*websocket.Conn
}

// Create a new PlayerConnections object
func NewPlayerConnections() *PlayerConnections {
	return &PlayerConnections{
		players: make(map[string]*websocket.Conn),
	}
}

func (pc *PlayerConnections) PingPlayerConnection(playerID string) {
	conn, exists := pc.GetPlayerConnection(playerID)
	if !exists {
		return
	}

	// Set up a ping/pong mechanism
	conn.SetPongHandler(func(appData string) error {
		// Pong received, keep the connection active
		return nil
	})

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Check if the connection is still valid before sending a ping
			if _, exists := pc.GetPlayerConnection(playerID); !exists {
				log.Println("Connection no longer exists, stopping ping.")
				return
			}
			
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Println("Ping failed, removing player:", playerID)
				pc.RemovePlayer(playerID)
				return
			}
		}
	}()
}

func (pc *PlayerConnections) AddPlayer(playerID string, conn *websocket.Conn) {
	fmt.Println("Starting AddPlayer for", playerID)
	pc.mu.Lock()
	defer pc.mu.Unlock()
	fmt.Println("Acquired lock for", playerID)
	pc.players[playerID] = conn
	fmt.Println("Added", playerID, "to players map")
	fmt.Println("About to call broadcastPlayerList for", playerID)
	go pc.broadcastPlayerList()  // Call broadcastPlayerList in a new goroutine
	fmt.Println("Finished AddPlayer for", playerID)
}

func (pc *PlayerConnections) RemovePlayer(playerID string) {
	fmt.Println("removing player", playerID)
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.players, playerID)
	pc.broadcastPlayerList()
}

func (pc *PlayerConnections) GetPlayerList() []string {
	fmt.Println("Entering GetPlayerList")
	pc.mu.RLock()
	fmt.Println("Acquired read lock in GetPlayerList")
	
	// Create a copy of the players map
	playersCopy := make(map[string]struct{})
	for playerID := range pc.players {
		playersCopy[playerID] = struct{}{}
		fmt.Printf("Copied player: %s\n", playerID)
	}
	pc.mu.RUnlock()
	fmt.Println("Released read lock in GetPlayerList")

	playerList := make([]string, 0, len(playersCopy))
	for playerID := range playersCopy {
		playerList = append(playerList, playerID)
		fmt.Printf("Added to list: %s\n", playerID)
	}
	fmt.Printf("GetPlayerList returning: %v\n", playerList)
	return playerList
}

func (pc *PlayerConnections) broadcastPlayerList() {
	fmt.Println("Starting to broadcast player list")
	time.Sleep(100 * time.Millisecond)  // Add a small delay
	playerList := pc.GetPlayerList()
	fmt.Printf("Player list: %v\n", playerList)
	
	message := map[string]interface{}{
		"type":    "playerList",
		"players": playerList,
	}
	
	fmt.Printf("Attempting to marshal message: %+v\n", message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshaling player list: %v\n", err)
		return
	}
	fmt.Println("JSON message created successfully")
	
	fmt.Printf("Attempting to broadcast message: %s\n", string(jsonMessage))
	err = pc.BroadcastMessage("playerList", playerList)
	if err != nil {
		fmt.Printf("Error broadcasting message: %v\n", err)
		return
	}
	
	fmt.Println("Player list broadcasted successfully")
}

func (pc *PlayerConnections) GetPlayerConnection(playerID string) (*websocket.Conn, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	conn, exists := pc.players[playerID]
	fmt.Printf("%s Exists? %t \n", playerID, exists);
	return conn, exists
}

// BroadcastMessage sends a message to all connected players
func (pc *PlayerConnections) BroadcastMessage(messageType string, message interface{}) error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// Create the JSON object
	jsonMessage := struct {
		Type    string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type:    messageType,
		Data: message,
	}

	// Marshal the JSON object
	jsonData, err := json.Marshal(jsonMessage)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	for playerID, conn := range pc.players {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Printf("Error broadcasting message to player %s: %v\n", playerID, err)
			return err
		}
	}
	return nil
}

// ShuffleAndDealCards shuffles the deck and deals cards to players
func (g *Game) ShuffleAndDealCards() {
	// Create a deck of 52 cards
	deck := createDeck()

	// Shuffle the deck
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	// Deal cards to players
	numPlayers := len(g.Players)
	for i, card := range deck {
		playerIndex := i % numPlayers
		//playerID := g.Players[playerIndex]
		
		// Find the player in the GameState and add the card to their hand
		switch playerIndex {
		case 0:
			g.State.Player1.Hand = append(g.State.Player1.Hand, card)
		case 1:
			g.State.Player2.Hand = append(g.State.Player2.Hand, card)
		case 2:
			g.State.Player3.Hand = append(g.State.Player3.Hand, card)
		case 3:
			g.State.Player4.Hand = append(g.State.Player4.Hand, card)
		}
	}
}

// createDeck generates a standard 52-card deck
func createDeck() []Card {
	deck := make([]Card, 0, 52)
	suits := []Suit{Hearts, Diamonds, Clubs, Spades}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}

	return deck
}
