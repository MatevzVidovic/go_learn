# 🎯 Complete Codebase Guide for Beginners

## 🏗️ Architecture Overview

Your online store runs **3 separate servers** working together:

| Server | Purpose | Port | Technology |
|--------|---------|------|------------|
| **Go API** | Handles HTTP requests, business logic | 8080 | Go + Gin framework |
| **MariaDB** | Stores all data (users, products, orders) | 3306 | SQL Database |
| **MQTT Broker** | Handles real-time messaging between services | 1883 | Eclipse Mosquitto |

## 📁 Code Structure Explained

### **What Each Folder Does**

```
online-store/
├── cmd/server/          # 🚀 Application entry point
├── internal/config/     # ⚙️  Environment variables & settings
├── internal/database/   # 🗄️  Database connection & table creation
├── internal/models/     # 📊 Data structures (User, Product, Order)
├── internal/handlers/   # 🌐 HTTP request processors
├── internal/middleware/ # 🛡️  Security & authentication
├── internal/mqtt/       # 📡 Real-time messaging
├── internal/services/   # 🧠 Business logic
└── docker-compose.yml   # 🐳 Runs all 3 servers together
```

### **The Layer Pattern (Very Important!)**

```
HTTP Request → Handler → Service → Database
                  ↓
               MQTT Event
```

Each layer has a specific job:
- **Handlers**: Receive HTTP requests, validate input
- **Services**: Contain business logic, make decisions
- **Database**: Store and retrieve data
- **MQTT**: Send events to other parts of the system

## 🌐 All Available Endpoints

### **Public Endpoints** (No login required)
```http
GET  /health                    # Check if server is running
GET  /api/products             # View all products
GET  /api/products/:id         # View specific product
POST /api/register             # Create new user account
POST /api/login                # Login and get access token
```

### **Protected Endpoints** (JWT token required)
```http
POST /api/products             # Create new product (admin)
PUT  /api/products/:id         # Update product (admin)
POST /api/orders               # Place an order
GET  /api/orders               # View your orders
GET  /api/orders/:id           # View specific order
```

### **How Authentication Works**
1. Client sends email/password to `/api/login`
2. Server creates a JWT token (like a temporary pass)
3. Client includes token in `Authorization: Bearer <token>` header
4. Middleware checks token before allowing access to protected endpoints

## 🔄 Request Flow Examples

### **Example 1: User Registration**
```
1. POST /api/register {"email": "user@test.com", "password": "pass123"}
2. AuthHandler.Register() validates the data
3. AuthService.Register() hashes password with bcrypt
4. Save user to database
5. Publish MQTT event "user/registered"
6. Return user info (without password) to client
```

### **Example 2: Creating an Order**
```
1. POST /api/orders {"product_id": 1, "quantity": 2}
2. Auth middleware checks JWT token
3. OrderHandler.CreateOrder() gets user ID from token
4. OrderService.CreateOrder() starts database transaction
5. Check if product has enough stock
6. Create order record, update product stock
7. Commit transaction (all-or-nothing)
8. Publish MQTT event "order/created"
9. Return order details to client
```

## 📡 MQTT Events Explained

**MQTT is like a message bus** - when something happens, we broadcast it so other parts can react.

### **Events We Publish** (Our app tells others)
```
user/registered     → "New user signed up"
user/login          → "User logged in"
product/created     → "New product added"
product/updated     → "Product details changed"
order/created       → "New order placed"
order/status_changed → "Order status updated"
inventory/low_stock → "Product running low"
```

### **Events We Listen For** (Others tell our app)
```
payment/confirmed   → "Payment was successful"
inventory/update    → "Stock levels changed"
```

### **Why Use MQTT?**
- **Decoupling**: Services don't need to know about each other directly
- **Real-time**: Instant updates across the system
- **Scalability**: Easy to add new services that react to events
- **Reliability**: Messages are queued if services are offline

## 🧠 Hard Concepts Explained Simply

### **1. Interfaces - The Go Way**
```go
type UserService interface {
    CreateUser(email string) error
}
```
**Think of interfaces as contracts.** Any struct that has a `CreateUser` method automatically implements this interface. This makes testing easier because you can create fake implementations.

### **2. Dependency Injection**
```go
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}
```
**Instead of creating dependencies inside functions, we pass them in.** This makes code more flexible and testable.

### **3. Error Handling Pattern**
```go
user, err := s.authService.Register(req)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```
**Go doesn't have exceptions.** Every function that can fail returns an error as the last value. You must check it explicitly.

### **4. Context Pattern**
```go
c.Set("user_id", userID)           // Middleware sets data
userID := c.Get("user_id").(int)   // Handler gets data
```
**Context carries data through the request lifecycle.** Middleware puts user info in context, handlers retrieve it.

### **5. Database Transactions**
```go
tx, err := s.db.Begin()            // Start transaction
defer tx.Rollback()                // Auto-rollback if something fails
// ... do database operations ...
tx.Commit()                        // Save all changes
```
**Transactions ensure data consistency.** Either all operations succeed, or none do.

## ⚠️ Important Warnings & Gotchas

### **Security Warnings**
- **JWT Secret**: Change the default secret in production!
- **Password Storage**: Never store plain text passwords (we use bcrypt)
- **MQTT Authentication**: Currently allows anonymous connections (development only)
- **Input Validation**: Always validate user input to prevent injection attacks

### **Database Warnings**
- **Connection Pooling**: We set max 25 connections - adjust based on load
- **Migration**: Currently recreates tables on startup - use proper migrations in production
- **Indexes**: Add database indexes for better performance on large datasets

### **MQTT Warnings**
- **QoS Levels**: We use QoS 1 (at least once delivery) - messages might duplicate
- **Message Size**: Keep MQTT messages small for better performance
- **Topic Design**: Use hierarchical topics (user/123/orders vs userorders123)

### **Go-Specific Gotchas**
- **Nil Pointers**: Check for nil before using pointers
- **Goroutine Leaks**: Always clean up goroutines (we handle this with defer)
- **Interface Assertion**: `userID := c.Get("user_id").(int)` panics if wrong type

### **Docker Warnings**
- **Data Persistence**: Database data is stored in Docker volumes
- **Port Conflicts**: Make sure ports 8080, 3306, 1883 are available
- **Health Checks**: Wait for health checks to pass before testing

## 🚀 Development Workflow

### **Making Changes**
1. **Modify Go code**
2. **Rebuild container**: `docker-compose up -d --build api`
3. **Test changes**: Use curl commands or Postman

### **Debugging**
```bash
# View logs
docker-compose logs -f api
docker-compose logs -f mariadb
docker-compose logs -f mqtt

# Connect to database
docker-compose exec mariadb mysql -u storeuser -pstorepass onlinestore

# Test MQTT
mosquitto_sub -h localhost -t "+"  # Listen to all topics
mosquitto_pub -h localhost -t "test" -m "hello"  # Send test message
```

### **Common Issues**
- **Port already in use**: Stop other services using the same ports
- **Database connection failed**: Wait for MariaDB health check to pass
- **MQTT not connecting**: Check if mosquitto.conf is properly mounted

## 🎓 Learning Path

### **Beginner Level**
1. Understand the request flow diagrams
2. Read through handlers to see how HTTP works
3. Look at models to understand data structures
4. Test endpoints with curl

### **Intermediate Level**
1. Study the service layer business logic
2. Understand JWT authentication flow
3. Learn about database transactions
4. Experiment with MQTT pub/sub

### **Advanced Level**
1. Add new features (shopping cart, payments)
2. Implement proper error handling
3. Add comprehensive tests
4. Optimize database queries

## 📈 Next Steps

Once you understand this codebase:
1. **Add features**: Shopping cart, email notifications, admin panel
2. **Improve security**: Rate limiting, input sanitization, HTTPS
3. **Scale up**: Load balancing, caching, microservices
4. **Monitor**: Logging, metrics, health checks

This codebase teaches you production-ready patterns that you'll see in real-world Go applications!