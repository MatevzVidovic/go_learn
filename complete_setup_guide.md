# 🚀 Complete Setup Guide

## Prerequisites

1. **Docker & Docker Compose**
   ```bash
   # Install Docker (Ubuntu)
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER
   
   # Install Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

2. **Go (for local development)**
   ```bash
   # Download and install Go
   wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

## Step-by-Step Setup

### 1. Create Project Structure

```bash
mkdir online-store
cd online-store

# Create directory structure
mkdir -p cmd/server
mkdir -p internal/{config,database,models,handlers,middleware,mqtt,services}
mkdir -p docs
```

### 2. Initialize Go Module

```bash
go mod init online-store

# Add dependencies
go get github.com/gin-gonic/gin
go get github.com/go-sql-driver/mysql
go get github.com/eclipse/paho.mqtt.golang
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto

# Clean up dependencies
go mod tidy
```

### 3. Create Source Files

Create all the files from the artifacts above in their respective directories:

```
online-store/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── database/connection.go
│   ├── models/ (user.go, product.go, order.go - all in one file)
│   ├── handlers/ (auth.go, products.go, orders.go - all in one file)
│   ├── middleware/auth.go
│   ├── mqtt/ (client.go, handlers.go - all in one file)
│   └── services/ (auth.go, products.go, orders.go - split into separate files)
├── docker-compose.yml
├── Dockerfile
├── mosquitto.conf
├── go.mod
└── README.md
```

### 4. Create Docker Configuration

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  mariadb:
    image: mariadb:10.6
    container_name: store_db
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: onlinestore
      MYSQL_USER: storeuser
      MYSQL_PASSWORD: storepass
    ports:
      - "3306:3306"
    volumes:
      - mariadb_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10

  mqtt:
    image: eclipse-mosquitto:2.0
    container_name: store_mqtt
    restart: always
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf
    command: mosquitto -c /mosquitto/config/mosquitto.conf

  api:
    build: .
    container_name: store_api
    restart: always
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: "storeuser:storepass@tcp(mariadb:3306)/onlinestore"
      MQTT_BROKER: "tcp://mqtt:1883"
      JWT_SECRET: "your-super-secret-jwt-key-change-this-in-production"
      PORT: "8080"
    depends_on:
      mariadb:
        condition: service_healthy
      mqtt:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      timeout: 10s
      retries: 5

volumes:
  mariadb_data:
```

Create `mosquitto.conf`:
```
allow_anonymous true
listener 1883
listener 9001
protocol websockets
log_type error
log_type warning
log_type notice
log_type information
log_dest stdout
persistence true
persistence_location /mosquitto/data/
```

### 5. Start the Application

```bash
# Build and start all services
docker-compose up -d

# Check if everything is running
docker-compose ps

# View logs
docker-compose logs -f
```

### 6. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# View products
curl http://localhost:8080/api/products
```

## Learning Path

### Beginner Concepts to Understand

1. **Structs** - Custom data types
   ```go
   type User struct {
       ID    int    `json:"id"`
       Email string `json:"email"`
   }
   ```

2. **Interfaces** - Contracts for behavior
   ```go
   type UserService interface {
       CreateUser(email string) error
   }
   ```

3. **Error Handling** - Explicit error checking
   ```go
   user, err := getUserByID(1)
   if err != nil {
       return err
   }
   ```

4. **Goroutines** - Concurrent execution
   ```go
   go func() {
       // This runs concurrently
   }()
   ```

5. **Channels** - Communication between goroutines
   ```go
   ch := make(chan string)
   go func() { ch <- "hello" }()
   msg := <-ch
   ```

### MQTT Concepts

1. **Publisher** - Sends messages to topics
2. **Subscriber** - Listens for messages on topics
3. **Topics** - Message categories (like "user/registered", "order/created")
4. **QoS** - Quality of Service (0=fire and forget, 1=at least once, 2=exactly once)
5. **Retained Messages** - Last message on topic is saved for new subscribers

### Key Go Patterns Used

1. **Constructor Functions**
   ```go
   func NewUserService(db *sql.DB) *UserService {
       return &UserService{db: db}
   }
   ```

2. **Error Wrapping**
   ```go
   if err != nil {
       return fmt.Errorf("failed to create user: %w", err)
   }
   ```

3. **Interface Composition**
   ```go
   type UserRepository interface {
       CreateUser(user User) error
       GetUser(id int) (*User, error)
   }
   ```

4. **Context Usage**
   ```go
   userID := c.Get("user_id").(int)
   ```

## Development Workflow

### Making Changes

1. **Modify code**
2. **Rebuild container**
   ```bash
   docker-compose up -d --build api
   ```

3. **Test changes**
   ```bash
   curl http://localhost:8080/health
   ```

### Adding New Features

1. **Add new model** (if needed)
2. **Create service methods**
3. **Add HTTP handlers**
4. **Update routes in main.go**
5. **Add MQTT events** (if applicable)
6. **Test with curl**

### Local Development (without Docker)

```bash
# Start database and MQTT with Docker
docker-compose up mariadb mqtt -d

# Set environment variables
export DATABASE_URL="storeuser:storepass@tcp(localhost:3306)/onlinestore"
export MQTT_BROKER="tcp://localhost:1883"
export JWT_SECRET="your-super-secret-jwt-key"
export PORT="8080"

# Run locally
go run cmd/server/main.go
```

## Common Issues and Solutions

### Database Connection Issues

```bash
# Check if MariaDB is running
docker-compose exec mariadb mysqladmin ping -h localhost

# Connect to database manually
docker-compose exec mariadb mysql -u storeuser -pstorepass onlinestore

# View tables
SHOW TABLES;
DESCRIBE users;
```

### MQTT Connection Issues

```bash
# Test MQTT broker
docker-compose exec mqtt mosquitto_pub -t "test" -m "hello"
docker-compose exec mqtt mosquitto_sub -t "test"

# Check MQTT logs
docker-compose logs mqtt
```

### Go Build Issues

```bash
# Clean module cache
go clean -modcache

# Rebuild dependencies
go mod tidy
go mod download

# Check for syntax errors
go vet ./...
```

## Next Steps - Adding Features

### 1. Shopping Cart

Add cart functionality:
```go
type Cart struct {
    ID     int `json:"id"`
    UserID int `json:"user_id"`
    Items  []CartItem `json:"items"`
}

type CartItem struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}
```

### 2. Payment Integration

Add Stripe payment:
```go
// In services/payments.go
func (s *PaymentService) CreatePaymentIntent(amount int) (*stripe.PaymentIntent, error) {
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(int64(amount)),
        Currency: stripe.String(string(stripe.CurrencyUSD)),
    }
    return paymentintent.New(params)
}
```

### 3. Email Notifications

Add email service:
```go
func (s *EmailService) SendOrderConfirmation(email string, order *models.Order) error {
    // Use SMTP or email service like SendGrid
}
```

### 4. WebSocket Support

Add real-time updates:
```go
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    // Handle WebSocket messages
}
```

### 5. Admin Dashboard

Add admin endpoints:
```go
// Admin-only routes
admin := api.Group("/admin")
admin.Use(middleware.AdminRequired())
{
    admin.GET("/users", adminHandler.GetUsers)
    admin.GET("/orders", adminHandler.GetAllOrders)
    admin.PUT("/orders/:id/status", adminHandler.UpdateOrderStatus)
}
```

## Production Considerations

### Security

1. **Use strong JWT secrets**
2. **Enable HTTPS**
3. **Add rate limiting**
4. **Validate all inputs**
5. **Use MQTT authentication**

### Performance

1. **Add database indexes**
2. **Use connection pooling**
3. **Implement caching (Redis)**
4. **Add monitoring**

### Deployment

1. **Use Docker Swarm or Kubernetes**
2. **Set up CI/CD pipeline**
3. **Configure load balancer**
4. **Set up logging and monitoring**

## Useful Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [MQTT Essentials](https://www.hivemq.com/mqtt-essentials/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [JWT.io](https://jwt.io/) - JWT debugger

## Summary

This project teaches you:

✅ **Go fundamentals** - structs, interfaces, error handling  
✅ **HTTP APIs** - REST endpoints with Gin framework  
✅ **Database operations** - SQL queries with Go  
✅ **MQTT messaging** - Event-driven architecture  
✅ **Authentication** - JWT tokens and middleware  
✅ **Docker** - Containerization and orchestration  
✅ **Testing** - Unit and integration tests  

The combination of HTTP and MQTT makes this a modern, scalable application architecture that you'll see in real-world systems!

Happy coding! 🎉