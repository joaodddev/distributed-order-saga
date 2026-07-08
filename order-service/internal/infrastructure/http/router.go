package http

import (
	"github.com/gin-gonic/gin"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/http/handler"
)

func NewRouter(orderHandler *handler.OrderHandler) *gin.Engine {
	router := gin.Default()

	orders := router.Group("/orders")
	{
		orders.POST("", orderHandler.Create)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	return router
}
