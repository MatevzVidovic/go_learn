# 🛍️ Go MQTT Online Store - Beginner Project

A complete online store backend built with Go, MariaDB, MQTT, and Docker. Perfect for learning Go and MQTT!

## 📁 Project Structure

```
online-store/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   ├── connection.go          # Database connection
│   │   └── migrations.sql         # Database schema
│   ├── models/
│   │   ├── user.go               # User data structures
│   │   ├── product.go            # Product data structures
│   │   └── order.go              # Order data structures
│   ├── handlers/
│   │   ├── auth.go               # Authentication endpoints
│   │   ├── products.go           # Product endpoints
│   │   └── orders.go             # Order endpoints
│   ├── middleware/
│   │   └── auth.go               # JWT authentication middleware
│   ├── mqtt/
│   │   ├── client.go             # MQTT client setup
│   │   └── handlers.go           # MQTT message handlers
│   └── services/
│       ├── auth.go               # Authentication business logic
│       ├── products.go           # Product business logic
│       └── orders.go             # Order business logic
├── docs/
│   └── swagger.yaml              # API documentation
├── docker-compose.yml            # Docker services setup
├── Dockerfile                    # Go application container
├── go.mod                        # Go modules file
├── go.sum                        # Go modules checksum
└── README.md                     # You are here!
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### 1. Clone and Run
```bash
# Clone the project
git clone <your-repo>
cd online-store

# Start all services with Docker
docker-compose up -d

# The API will be available at http://localhost:8080
# Swagger docs at http://localhost:8080/swagger/
```

### 2. Test the API
```bash
# Register a new user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 📚 Learning Guide

### Step 1: Understanding the Structure
- `cmd/server/main.go` - This is where your application starts
- `internal/` - All your application code (internal means other projects can't import it)
- `models/` - Data structures (like User, Product)
- `handlers/` - HTTP endpoint functions
- `services/` - Business logic (the "brain" of your app)

### Step 2: Key Go Concepts You'll Learn
- **Structs** - Custom data types (like classes in other languages)
- **Interfaces** - Contracts that types must fulfill
- **Goroutines** - Lightweight threads for concurrent programming
- **Channels** - Communication between goroutines
- **Error handling** - Go's explicit error handling pattern

### Step 3: MQTT Integration
- **Publisher** - Sends messages when things happen (new order, user registered)
- **Subscriber** - Listens for messages and reacts
- **Topics** - Message categories (like "user/registered", "order/created")

## 🔧 Configuration

Environment variables (set in docker-compose.yml):
```env
DB_HOST=mariadb
DB_PORT=3306
DB_USER=storeuser
DB_PASSWORD=storepass
DB_NAME=onlinestore
MQTT_BROKER=tcp://mqtt:1883
JWT_SECRET=your-secret-key
```

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table
```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_cents INT NOT NULL,
    stock_quantity INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Orders Table
```sql
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    total_cents INT NOT NULL,
    status ENUM('pending', 'paid', 'shipped', 'delivered') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## 🔄 MQTT Topics

### Published Events (what our app sends)
- `user/registered` - When a new user signs up
- `user/login` - When a user logs in
- `product/created` - When a new product is added
- `product/updated` - When product details change
- `order/created` - When a new order is placed
- `order/status_changed` - When order status updates

### Subscribed Events (what our app listens for)
- `inventory/low_stock` - Alert when product stock is low
- `payment/confirmed` - When payment is processed
- `notifications/send` - Trigger to send emails/SMS

## 🛠️ API Endpoints

### Authentication
- `POST /api/register` - Create new user account
- `POST /api/login` - Login and get JWT token

### Products
- `GET /api/products` - List all products
- `POST /api/products` - Create new product (admin only)
- `GET /api/products/:id` - Get specific product
- `PUT /api/products/:id` - Update product (admin only)

### Orders
- `POST /api/orders` - Create new order
- `GET /api/orders` - Get user's orders
- `GET /api/orders/:id` - Get specific order

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires Docker)
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

## 📖 Learning Resources

### Go Basics
- [Go Tour](https://tour.golang.org/) - Interactive Go tutorial
- [Effective Go](https://golang.org/doc/effective_go.html) - Go best practices

### MQTT
- [MQTT Essentials](https://www.hivemq.com/mqtt-essentials/) - Complete MQTT guide
- [Eclipse Paho Go Client](https://github.com/eclipse/paho.mqtt.golang) - MQTT client we use

### Advanced Topics
- [Go Concurrency Patterns](https://blog.golang.org/pipelines) - Goroutines and channels
- [REST API Design](https://restfulapi.net/) - API best practices

## 🚀 Next Steps

1. **Add more features**: Shopping cart, payment integration, email notifications
2. **Improve security**: Rate limiting, input validation, HTTPS
3. **Scale up**: Load balancing, caching, database optimization
4. **Monitor**: Logging, metrics, health checks

## 🤝 Contributing

This is a learning project! Feel free to:
- Add new features
- Improve documentation
- Fix bugs
- Share your learning experience

## 📝 Common Beginner Questions

**Q: What's the difference between `internal/` and `pkg/`?**
A: `internal/` code can only be imported by this project. `pkg/` can be imported by other projects.

**Q: Why use interfaces?**
A: They make testing easier and code more flexible. You can swap implementations without changing other code.

**Q: When should I use goroutines?**
A: For I/O operations (database, HTTP calls, MQTT) that can run concurrently without blocking.

**Q: How does MQTT help?**
A: It decouples services. Instead of direct HTTP calls, services communicate through messages, making the system more flexible and scalable.