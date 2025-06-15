// internal/database/connection.go
// This file handles connecting to our MariaDB database

package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver (MariaDB is compatible)
)

// Connect creates a connection to the database
// It returns a *sql.DB which is a connection pool, not a single connection
func Connect(databaseURL string) (*sql.DB, error) {
	// Open creates a database connection pool
	// It doesn't actually connect yet - that happens on first query
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection by pinging the database
	// This is where we actually connect for the first time
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	// MaxOpenConns: maximum number of open connections to the database
	// MaxIdleConns: maximum number of connections in the idle connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Create tables if they don't exist
	// In a real app, you'd use migrations, but this is simpler for learning
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates all the database tables we need
// This runs every time the app starts, but only creates tables if they don't exist
func createTables(db *sql.DB) error {
	// SQL queries to create our tables
	// These are the same as in our README, but in Go string format
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price_cents INT NOT NULL,
			stock_quantity INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			product_id INT NOT NULL,
			quantity INT NOT NULL,
			total_cents INT NOT NULL,
			status ENUM('pending', 'paid', 'shipped', 'delivered') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
	}

	// Execute each CREATE TABLE query
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	// Insert some sample products if the products table is empty
	// This gives us data to work with when testing
	if err := insertSampleData(db); err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	return nil
}

// insertSampleData adds some example products to the database
// This only runs if the products table is empty
func insertSampleData(db *sql.DB) error {
	// Check if we already have products
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}

	// If we already have products, don't add more
	if count > 0 {
		return nil
	}

	// Sample products to insert
	products := []struct {
		name        string
		description string
		priceCents  int
		stock       int
	}{
		{"Go Programming Book", "Learn Go programming from scratch", 2999, 50},
		{"MQTT Sensor Kit", "IoT sensor kit with MQTT support", 4999, 25},
		{"Docker T-Shirt", "Comfortable cotton t-shirt with Docker logo", 1999, 100},
		{"Wireless Mouse", "Ergonomic wireless mouse for developers", 3499, 75},
	}

	// Insert each sample product
	for _, product := range products {
		_, err := db.Exec(
			"INSERT INTO products (name, description, price_cents, stock_quantity) VALUES (?, ?, ?, ?)",
			product.name, product.description, product.priceCents, product.stock,
		)
		if err != nil {
			return err
		}
	}

	return nil
}