// internal/handlers/auth.go
// This file contains HTTP handlers for authentication endpoints

package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"online-store/internal/models"
	"online-store/internal/services"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration requests
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegistration true "User registration data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration
	
	// Bind JSON request to struct and validate
	// Gin will automatically check the binding rules we defined in the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to register the user
	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the created user (without password)
	c.JSON(http.StatusCreated, user)
}

// Login handles user login requests
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to login the user
	token, user, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return the token and user info
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
