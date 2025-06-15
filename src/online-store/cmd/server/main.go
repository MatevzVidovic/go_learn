// cmd/server/main.go
// This is the entry point of our application - where everything starts!

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"online-store/internal/config"
	"online-store/internal/database"
	"online-store/internal/handlers"
	"online-store/internal/middleware"
	"online-store/internal/mqtt"
	"online-store/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment variables
	// This is where we get database connection info, MQTT settings, etc.
	cfg := config.Load()

	// Connect to the database (MariaDB)
	// This creates a connection pool that our app will use
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close() // Make sure we close the connection when the app shuts down

	// Set up MQTT client for publishing and subscribing to messages
	// MQTT helps different parts of our system communicate
	mqttClient, err := mqtt.NewClient(cfg.MQTTBroker)
	if err != nil {
		log.Fatal("Failed to connect to MQTT broker:", err)
	}
	defer mqttClient.Disconnect(250) // Clean disconnect when shutting down

	// Create service layer - this is where our business logic lives
	// Services handle the "what" and "how" of our application
	authService := services.NewAuthService(db, mqttClient)
	productService := services.NewProductService(db, mqttClient)
	orderService := services.NewOrderService(db, mqttClient)

	// Create HTTP handlers - these handle incoming web requests
	// Handlers are like receptionists that greet requests and hand them off
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Set up MQTT message handlers
	// These listen for MQTT messages and do something when they arrive
	mqttHandlers := mqtt.NewHandlers(productService, orderService)
	mqttHandlers.Subscribe(mqttClient)

	// Create Gin router (Gin is a web framework for Go)
	// Think of this as the traffic director for web requests
	router := gin.Default()

	// Add middleware - code that runs before every request
	// CORS allows web browsers to make requests to our API
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define API routes - these are the URLs our app responds to
	api := router.Group("/api")
	{
		// Authentication routes - no middleware needed, anyone can access
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Product routes - some need authentication, some don't
		api.GET("/products", productHandler.GetProducts)    // Anyone can view products
		api.GET("/products/:id", productHandler.GetProduct) // Anyone can view a product

		// Protected routes - need to be logged in (JWT token required)
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(cfg.JWTSecret)) // Check if user is logged in
		{
			// Only logged-in users can create products, orders, etc.
			protected.POST("/products", productHandler.CreateProduct)
			protected.PUT("/products/:id", productHandler.UpdateProduct)
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.GetUserOrders)
			protected.GET("/orders/:id", orderHandler.GetOrder)
		}
	}

	// Health check endpoint - useful for monitoring if the app is running
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// Start the HTTP server in a goroutine (concurrent execution)
	// This means the server runs in the background while we wait for shutdown signals
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Set up graceful shutdown
	// This catches Ctrl+C and other shutdown signals to close cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait for shutdown signal

	log.Println("Shutting down server...")
	// App will automatically clean up database and MQTT connections due to defer statements above
}
