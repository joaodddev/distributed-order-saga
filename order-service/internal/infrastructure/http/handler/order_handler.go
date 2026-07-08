package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/input"
)

type OrderHandler struct {
	createOrder input.CreateOrderUseCase
}

func NewOrderHandler(createOrder input.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{createOrder: createOrder}
}

type createOrderItemRequest struct {
	ProductID string  `json:"productId" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unitPrice" binding:"required,gt=0"`
}

type createOrderRequest struct {
	CustomerID string                   `json:"customerId" binding:"required"`
	Items      []createOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customerId"})
		return
	}

	items := make([]input.CreateOrderItemInput, len(req.Items))
	for i, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productId"})
			return
		}
		items[i] = input.CreateOrderItemInput{
			ProductID: productID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	output, err := h.createOrder.Execute(c.Request.Context(), input.CreateOrderInput{
		CustomerID: customerID,
		Items:      items,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"orderId": output.OrderID,
		"status":  output.Status,
	})
}
