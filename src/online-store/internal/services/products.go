
// internal/services/products.go
// This file contains product-related business logic

// ProductService handles product operations
type ProductService struct {
	db         *sql.DB
	mqttClient *mqtt.Client
}

// NewProductService creates a new product service
func NewProductService(db *sql.DB, mqttClient *mqtt.Client) *ProductService {
	return &ProductService{
		db:         db,
		mqttClient: mqttClient,
	}
}

// GetProducts returns all products
func (s *ProductService) GetProducts() ([]models.Product, error) {
	rows, err := s.db.Query(
		"SELECT id, name, description, price_cents, stock_quantity, created_at FROM products ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close() // Always close rows when done

	var products []models.Product
	
	// Iterate through all rows
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceCents,
			&product.StockQuantity,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// GetProduct returns a single product by ID
func (s *ProductService) GetProduct(id int) (*models.Product, error) {
	var product models.Product
	err := s.db.QueryRow(
		"SELECT id, name, description, price_cents, stock_quantity, created_at FROM products WHERE id = ?",
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.PriceCents,
		&product.StockQuantity,
		&product.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req models.ProductRequest) (*models.Product, error) {
	result, err := s.db.Exec(
		"INSERT INTO products (name, description, price_cents, stock_quantity) VALUES (?, ?, ?, ?)",
		req.Name, req.Description, req.PriceCents, req.StockQuantity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get product ID: %w", err)
	}

	// Get the created product
	product, err := s.GetProduct(int(productID))
	if err != nil {
		return nil, err
	}

	// Publish MQTT event
	event := models.ProductCreatedEvent{
		ProductID: product.ID,
		Name:      product.Name,
		Timestamp: time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("product/created", event); err != nil {
		fmt.Printf("Failed to publish product created event: %v", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id int, req models.ProductRequest) (*models.Product, error) {
	_, err := s.db.Exec(
		"UPDATE products SET name = ?, description = ?, price_cents = ?, stock_quantity = ? WHERE id = ?",
		req.Name, req.Description, req.PriceCents, req.StockQuantity, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Get the updated product
	product, err := s.GetProduct(id)
	if err != nil {
		return nil, err
	}

	// Publish MQTT event
	event := struct {
		ProductID int   `json:"product_id"`
		Name      string `json:"name"`
		Timestamp int64  `json:"timestamp"`
	}{
		ProductID: product.ID,
		Name:      product.Name,
		Timestamp: time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("product/updated", event); err != nil {
		fmt.Printf("Failed to publish product updated event: %v", err)
	}

	return product, nil
}

// UpdateStock updates the stock quantity for a product
// This method is called by MQTT handlers
func (s *ProductService) UpdateStock(productID, newStock int) error {
	_, err := s.db.Exec(
		"UPDATE products SET stock_quantity = ? WHERE id = ?",
		newStock, productID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Check if stock is low (less than 10 items)
	if newStock < 10 {
		product, err := s.GetProduct(productID)
		if err != nil {
			return err
		}

		// Send low stock alert
		alert := models.LowStockAlert{
			ProductID:    productID,
			ProductName:  product.Name,
			CurrentStock: newStock,
			ReorderLevel: 10,
			Timestamp:    time.Now().Unix(),
		}
		
		if err := s.mqttClient.Publish("inventory/low_stock", alert); err != nil {
			fmt.Printf("Failed to publish low stock alert: %v", err)
		}
	}

	return nil
}