// internal/mqtt/client.go
// This file sets up our MQTT client for publishing and subscribing to messages

package mqtt

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client wraps the MQTT client with our custom methods
type Client struct {
	client MQTT.Client
}

// NewClient creates a new MQTT client and connects to the broker
func NewClient(brokerURL string) (*Client, error) {
	// Generate a random client ID
	// Each MQTT client needs a unique ID
	clientID := generateClientID()

	// Set up MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerURL)   // Where to connect
	opts.SetClientID(clientID)  // Our unique ID
	opts.SetCleanSession(true)  // Start fresh each time
	opts.SetAutoReconnect(true) // Reconnect if connection drops
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetKeepAlive(30 * time.Second)

	// Set up connection handlers
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	opts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Println("MQTT client connected")
	})

	// Create the client
	client := MQTT.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return &Client{client: client}, nil
}

// Publish sends a message to an MQTT topic
// This is how we tell other parts of the system that something happened
func (c *Client) Publish(topic string, payload interface{}) error {
	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Publish the message
	// QoS 1 means "at least once delivery" - the message will be delivered at least once
	// false means "not retained" - the broker won't save this message for future subscribers
	token := c.client.Publish(topic, 1, false, jsonData)

	// Wait for the publish to complete
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}

	log.Printf("Published message to topic %s: %s", topic, string(jsonData))
	return nil
}

// Subscribe listens for messages on an MQTT topic
// When a message arrives, it calls the provided handler function
func (c *Client) Subscribe(topic string, handler MQTT.MessageHandler) error {
	// Subscribe to the topic
	// QoS 1 means we want reliable delivery
	token := c.client.Subscribe(topic, 1, handler)

	// Wait for the subscription to complete
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Disconnect closes the MQTT connection
func (c *Client) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
	log.Println("MQTT client disconnected")
}

// generateClientID creates a random client ID for MQTT
func generateClientID() string {
	// Create a random 8-byte array
	bytes := make([]byte, 8)
	rand.Read(bytes)

	// Convert to hex string and add prefix
	return fmt.Sprintf("store-client-%x", bytes)
}
