// internal/mqtt/handlers.go
// This file contains MQTT message handlers - functions that run when messages arrive

package mqtt

import (
	"encoding/json"
	"log"
	"online-store/internal/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Handlers manages all MQTT message subscriptions and handlers
type Handlers struct {
	productService ProductService // Interface for product operations
	orderService   OrderService   // Interface for order operations
}

// ProductService interface defines what product operations we need
// Using interfaces makes testing easier and code more flexible
type ProductService interface {
	UpdateStock(productID, newStock int) error
	GetProduct(id int) (*models.Product, error)
}

// OrderService interface defines what order operations we need
type OrderService interface {
	UpdateOrderStatus(orderID int, status string) error
}

// NewHandlers creates a new MQTT handlers manager
func NewHandlers(productService ProductService, orderService OrderService) *Handlers {
	return &Handlers{
		productService: productService,
		orderService:   orderService,
	}
}

// Subscribe sets up all our MQTT subscriptions
// This is where we tell MQTT what topics we want to listen to
func (h *Handlers) Subscribe(client *Client) {
	// Subscribe to inventory updates
	client.Subscribe("inventory/update", h.handleInventoryUpdate)

	// Subscribe to payment confirmations
	client.Subscribe("payment/confirmed", h.handlePaymentConfirmed)

	// Subscribe to stock alerts
	client.Subscribe("inventory/low_stock", h.handleLowStockAlert)

	log.Println("All MQTT subscriptions set up")
}

// handleInventoryUpdate processes inventory update messages
func (h *Handlers) handleInventoryUpdate(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received inventory update: %s", string(msg.Payload()))

	// Parse the message
	var update struct {
		ProductID int `json:"product_id"`
		NewStock  int `json:"new_stock"`
	}

	if err := json.Unmarshal(msg.Payload(), &update); err != nil {
		log.Printf("Failed to parse inventory update: %v", err)
		return
	}

	// Update the product stock
	if err := h.productService.UpdateStock(update.ProductID, update.NewStock); err != nil {
		log.Printf("Failed to update product stock: %v", err)
		return
	}

	log.Printf("Updated stock for product %d to %d", update.ProductID, update.NewStock)
}

// handlePaymentConfirmed processes payment confirmation messages
func (h *Handlers) handlePaymentConfirmed(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received payment confirmation: %s", string(msg.Payload()))

	// Parse the message
	var payment struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"`
	}

	if err := json.Unmarshal(msg.Payload(), &payment); err != nil {
		log.Printf("Failed to parse payment confirmation: %v", err)
		return
	}

	// Update the order status
	if err := h.orderService.UpdateOrderStatus(payment.OrderID, "paid"); err != nil {
		log.Printf("Failed to update order status: %v", err)
		return
	}

	log.Printf("Updated order %d status to paid", payment.OrderID)
}

// handleLowStockAlert processes low stock alert messages
func (h *Handlers) handleLowStockAlert(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received low stock alert: %s", string(msg.Payload()))

	// In a real application, you might:
	// 1. Send an email to the inventory manager
	// 2. Automatically reorder products
	// 3. Update a dashboard
	// 4. Log to a monitoring system

	// For this example, we'll just log it
	var alert models.LowStockAlert
	if err := json.Unmarshal(msg.Payload(), &alert); err != nil {
		log.Printf("Failed to parse low stock alert: %v", err)
		return
	}

	log.Printf("LOW STOCK ALERT: Product %s (ID: %d) has only %d items left!",
		alert.ProductName, alert.ProductID, alert.CurrentStock)
}
