// internal/auth/jwt.go
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWTs
var jwtKey = []byte("my_secret_key")

// Claims structure for the JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Generate a new JWT token for a given username
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Token expires in 15 minutes

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT checks the JWT token and extracts claims
func ValidateJWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get the Authorization header
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        // Strip out the "Bearer " part if present
        const bearerPrefix = "Bearer "
        if len(tokenString) > len(bearerPrefix) && tokenString[:len(bearerPrefix)] == bearerPrefix {
            tokenString = tokenString[len(bearerPrefix):] // Extract the token part
        } else {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        // Validate the token (your function)
        claims, err := ValidateJWT(tokenString)
        if err != nil {
            fmt.Printf("Token validation failed: %v\n", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        fmt.Println("claims", claims)
        // Optional: Add claims to context
        ctx := context.WithValue(r.Context(), claimsKey, claims)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}

func GetUsernameFromToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the username from the claims
		username, ok := claims["username"].(string)
		if !ok {
			return "", fmt.Errorf("username claim not found in token")
		}
		return username, nil
	}

	return "", fmt.Errorf("invalid token")
}


// Validate a JWT token and return the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Define a custom type for the context key
type contextKey string

// Create a constant for the claims key
const claimsKey contextKey = "claims"
