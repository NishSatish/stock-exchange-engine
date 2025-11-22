package api

import (
	"fmt"
	"net/http"
	"stock-exchange-simulator/api/dto"
	"stock-exchange-simulator/pkg/app"
	"stock-exchange-simulator/pkg/models"

	"github.com/gin-gonic/gin"
)

// PlaceOrderHandler creates a Gin handler function for placing orders.
func PlaceOrderHandler(appServices *app.AppServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.PlaceOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "some shi missing in yo input"})
			return
		}
		fmt.Print("ww", req)

		order := &models.Order{
			StockID:  req.StockID,
			Type:     req.Type,
			Quantity: req.Quantity,
			Price:    req.Price,
			Status:   "pending", // Default status for new orders
		}

		exchangeOrderId, err := appServices.Nexus.PlaceOrder(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "exchange_order_id": exchangeOrderId})
	}
}
