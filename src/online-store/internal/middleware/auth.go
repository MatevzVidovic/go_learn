// internal/middleware/auth.go
// This file contains JWT authentication middleware

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired is middleware that checks for valid JWT tokens
// Middleware is code that runs before your actual handler functions
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get the Authorization header
		// Format should be: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort() // Stop processing, don't call the next handler
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Return our secret key for validation
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is valid and get claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user information from token
			userID, ok := claims["user_id"].(float64) // JSON numbers are float64 in Go
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			email, ok := claims["email"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			// Store user information in the context so handlers can access it
			// This is how we pass data from middleware to handlers
			c.Set("user_id", int(userID))
			c.Set("user_email", email)

			// Continue to the next handler
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	})
}