package handlers

import (
	//"encoding/json"
	// "fmt"
	//"dealer-backend/internal/auth"
	"dealer-backend/internal/config"
	"dealer-backend/internal/services"
	"io"
	"net/http"
)


func StartHandler(w http.ResponseWriter, r *http.Request) {
    playerID := r.URL.Query().Get("playerID")
    services.StartMatchmaking(config.PlayerConnections, playerID)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        // Parse `gameID` and `playerID` from the query parameters
        gameID := r.URL.Query().Get("gameID")
        playerID := r.URL.Query().Get("playerID")

        if gameID == "" || playerID == "" {
            http.Error(w, "gameID and playerID are required", http.StatusBadRequest)
            return
        }

        // Read the move data (message body) from the request
        message, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read move data", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // Call handlePlayerMove with parsed gameID, playerID, and move data
        services.HandlePlayerMove(gameID, playerID, message)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Move made"))
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}





