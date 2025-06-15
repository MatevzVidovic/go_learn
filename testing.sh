# Testing and Usage Examples

# 1. Start the application
docker-compose up -d

# Wait for services to be healthy
docker-compose ps

# 2. Test user registration
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Expected response:
# {
#   "id": 1,
#   "email": "john@example.com",
#   "created_at": "2024-01-15T10:30:00Z"
# }

# 3. Test user login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Expected response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": 1,
#     "email": "john@example.com",
#     "created_at": "2024-01-15T10:30:00Z"
#   }
# }

# Save the token for authenticated requests
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 4. View all products (no authentication needed)
curl http://localhost:8080/api/products

# Expected response:
# [
#   {
#     "id": 1,
#     "name": "Go Programming Book",
#     "description": "Learn Go programming from scratch",
#     "price_cents": 2999,
#     "stock_quantity": 50,
#     "created_at": "2024-01-15T10:00:00Z"
#   },
#   ...
# ]

# 5. View a specific product
curl http://localhost:8080/api/products/1

# 6. Create a new product (authentication required)
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Raspberry Pi Kit",
    "description": "Complete Raspberry Pi starter kit",
    "price_cents": 7999,
    "stock_quantity": 30
  }'

# 7. Create an order (authentication required)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'

# Expected response:
# {
#   "id": 1,
#   "product_id": 1,
#   "product_name": "Go Programming Book",
#   "quantity": 2,
#   "total_cents": 5998,
#   "status": "pending",
#   "created_at": "2024-01-15T10:45:00Z"
# }

# 8. View your orders
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/orders

# 9. View a specific order
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/orders/1

# 10. Check application health
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "ok",
#   "timestamp": "2024-01-15T10:50:00Z"
# }

# MQTT Testing Examples

# 11. Listen to MQTT messages (install mosquitto-clients first)
# On Ubuntu/Debian: sudo apt-get install mosquitto-clients
# On macOS: brew install mosquitto

# Subscribe to all events
mosquitto_sub -h localhost -t "user/+"
mosquitto_sub -h localhost -t "product/+"
mosquitto_sub -h localhost -t "order/+"
mosquitto_sub -h localhost -t "inventory/+"

# 12. Simulate external MQTT messages

# Simulate payment confirmation
mosquitto_pub -h localhost -t "payment/confirmed" -m '{
  "order_id": 1,
  "status": "paid"
}'

# Simulate inventory update
mosquitto_pub -h localhost -t "inventory/update" -m '{
  "product_id": 1,
  "new_stock": 5
}'

# Development Commands

# 13. Run locally (without Docker)
export DATABASE_URL="storeuser:storepass@tcp(localhost:3306)/onlinestore"
export MQTT_BROKER="tcp://localhost:1883"
export JWT_SECRET="your-super-secret-jwt-key"
export PORT="8080"

go run cmd/server/main.go

# 14. Run tests
go test ./...

# 15. Build the application
go build -o store-api ./cmd/server

# 16. Check logs
docker-compose logs -f api
docker-compose logs -f mqtt
docker-compose logs -f mariadb

# 17. Stop everything
docker-compose down

# 18. Clean up (removes volumes too)
docker-compose down -v

# Troubleshooting

# Check if services are running
docker-compose ps

# Check service health
docker-compose exec api wget -q --spider http://localhost:8080/health
docker-compose exec mariadb mysqladmin ping -h localhost

# Connect to database directly
docker-compose exec mariadb mysql -u storeuser -pstorepass onlinestore

# Connect to MQTT broker
mosquitto_pub -h localhost -t "test" -m "hello world"
mosquitto_sub -h localhost -t "test"