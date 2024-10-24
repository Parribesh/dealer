// internal/handlers/user.go
package handlers

import (
	"dealer-backend/internal/auth"
	"fmt"
	"net/http"
)

// Handler for user login (JWT generation)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	// Generate JWT for the user
	token, err := auth.GenerateJWT(username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the JWT in the response
	w.Write([]byte(fmt.Sprintf("Token: %s", token)))
}

// Handler for protected route (JWT validation)
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Validate the token
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Allow access to the protected resource
	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Username)))
}
