// internal/models/user.go
// This file defines the User data structure

package models

import "time"

// User represents a user in our system
// In Go, we use structs to define data structures
type User struct {
	ID           int       `json:"id" db:"id"`                         // Database ID
	Email        string    `json:"email" db:"email"`                   // User's email address
	PasswordHash string    `json:"-" db:"password_hash"`               // Hashed password (json:"-" means don't include in JSON)
	CreatedAt    time.Time `json:"created_at" db:"created_at"`         // When the user was created
}

// UserRegistration represents the data needed to register a new user
// We separate this from User because we don't want to expose password hashes
type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`    // Email is required and must be valid email format
	Password string `json:"password" binding:"required,min=6"` // Password is required and must be at least 6 characters
}

// UserLogin represents login credentials
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is what we send back to the client (without sensitive data)
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to UserResponse (removes sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

