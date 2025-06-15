
// internal/handlers/products.go
// This file contains HTTP handlers for product endpoints

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts returns all products
// @Summary Get all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Router /api/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct returns a specific product
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Get ID from URL parameter
	id, err := getIDFromParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.ProductRequest true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
// @Summary Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.ProductRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := getIDFromParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
