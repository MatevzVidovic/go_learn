
// internal/handlers/orders.go
// This file contains HTTP handlers for order endpoints

// OrderHandler handles order HTTP requests
type OrderHandler struct {
	orderService *services.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.OrderRequest true "Order data"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetUserOrders returns all orders for the authenticated user
// @Summary Get user's orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.OrderResponse
// @Security BearerAuth
// @Router /api/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder returns a specific order for the authenticated user
// @Summary Get order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.OrderResponse
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orderID, err := getIDFromParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(orderID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Helper functions

// getIDFromParam extracts an integer ID from URL parameters
func getIDFromParam(c *gin.Context, param string) (int, error) {
	// strconv package is used to convert strings to other types
	idStr := c.Param(param)
	
	// Convert string to integer
	// ParseInt(string, base, bitSize) - base 10 = decimal, bitSize 0 = int
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return 0, err
	}
	
	return int(id), nil
}

// getUserIDFromContext extracts user ID from the Gin context
// This is set by our authentication middleware
func getUserIDFromContext(c *gin.Context) (int, error) {
	// Get user ID that was set by auth middleware
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	// Convert interface{} to int
	userID, ok := userIDInterface.(int)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	return userID, nil
}