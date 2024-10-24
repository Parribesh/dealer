package handlers

import (
	//"encoding/json"
	// "fmt"
	//"dealer-backend/internal/auth"
	"dealer-backend/internal/config"
	"dealer-backend/internal/services"
	"net/http"
)


func StartHandler(w http.ResponseWriter, r *http.Request) {
    playerID := r.URL.Query().Get("playerID")
    services.StartMatchmaking(config.PlayerConnections, playerID)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        // Handle POST request logic for making a move
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Move made"))
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
