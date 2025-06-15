// internal/models/product.go
// Product represents a product in our store

// Product represents an item in our online store
type Product struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	PriceCents    int       `json:"price_cents" db:"price_cents"`       // Price in cents (avoids floating point issues)
	StockQuantity int       `json:"stock_quantity" db:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ProductRequest represents data needed to create/update a product
type ProductRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PriceCents    int    `json:"price_cents" binding:"required,min=1"`     // Must be at least 1 cent
	StockQuantity int    `json:"stock_quantity" binding:"required,min=0"`  // Can't have negative stock
}

// PriceInDollars returns the price in dollars (for display purposes)
func (p *Product) PriceInDollars() float64 {
	return float64(p.PriceCents) / 100.0
}
