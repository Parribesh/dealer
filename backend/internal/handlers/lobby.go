package handlers

import (
	"crypto/rand"
	"dealer-backend/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

// GenerateRandomUsername creates a random username for a new player
func GenerateRandomUsername() string {
    adjectives := []string{"Brave", "Swift", "Clever", "Fierce", "Lucky"}
    animals := []string{"Lion", "Eagle", "Tiger", "Wolf", "Falcon"}

    randIndex1, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
    randIndex2, _ := rand.Int(rand.Reader, big.NewInt(int64(len(animals))))

    return fmt.Sprintf("%s%s", adjectives[randIndex1.Int64()], animals[randIndex2.Int64()])
}

func JoinHandler(w http.ResponseWriter, r *http.Request) {
    // Generate a random username
    username := r.URL.Query().Get("username") 
    if username == "" {
        username = GenerateRandomUsername() 
    }
    // username:= "random"    
    log.Printf("Username set to: %s", username)
    // Set the response header to indicate JSON content
    w.Header().Set("Content-Type", "application/json")
    
	token, err := auth.GenerateJWT(username)

	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Create a JSON response
    response := map[string]string{"token": token, "playerName": username}
    // Write the JSON response
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Send the response with a 200 OK status
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}


func LobbyStatusHandler(w http.ResponseWriter, r *http.Request){
	
}