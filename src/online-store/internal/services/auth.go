// internal/services/auth.go
// This file contains authentication business logic

package services

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"online-store/internal/models"
	"online-store/internal/mqtt"
)

// AuthService handles user authentication operations
type AuthService struct {
	db         *sql.DB      // Database connection
	mqttClient *mqtt.Client // MQTT client for publishing events
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, mqttClient *mqtt.Client) *AuthService {
	return &AuthService{
		db:         db,
		mqttClient: mqttClient,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req models.UserRegistration) (*models.UserResponse, error) {
	// Hash the password using bcrypt
	// bcrypt is a secure way to store passwords - it's slow and uses salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert the user into the database
	result, err := s.db.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		req.Email, string(hashedPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the ID of the newly created user
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Create user response
	userResponse := &models.UserResponse{
		ID:        int(userID),
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	// Publish MQTT event that a new user registered
	// This allows other parts of the system to react (send welcome email, etc.)
	event := models.UserRegisteredEvent{
		UserID:    int(userID),
		Email:     req.Email,
		Timestamp: time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("user/registered", event); err != nil {
		// Don't fail the registration if MQTT publish fails
		// Just log the error - the user was created successfully
		fmt.Printf("Failed to publish user registered event: %v", err)
	}

	return userResponse, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req models.UserLogin) (string, *models.UserResponse, error) {
	// Get user from database
	var user models.User
	err := s.db.QueryRow(
		"SELECT id, email, password_hash, created_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, fmt.Errorf("invalid email or password")
		}
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Create JWT token
	token, err := s.createJWTToken(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Publish MQTT event that user logged in
	event := struct {
		UserID    int   `json:"user_id"`
		Email     string `json:"email"`
		Timestamp int64  `json:"timestamp"`
	}{
		UserID:    user.ID,
		Email:     user.Email,
		Timestamp: time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("user/login", event); err != nil {
		fmt.Printf("Failed to publish user login event: %v", err)
	}

	userResponse := user.ToResponse()
	return token, &userResponse, nil
}

// createJWTToken creates a JWT token for a user
func (s *AuthService) createJWTToken(userID int, email string) (string, error) {
	// JWT claims - the data we put inside the token
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with our secret key
	// In production, use a strong random secret key
	return token.SignedString([]byte("your-super-secret-jwt-key-change-this-in-production"))
}
