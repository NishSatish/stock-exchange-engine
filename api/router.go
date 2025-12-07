package api

import (
	"stock-exchange-simulator/pkg/app"

	"github.com/gin-gonic/gin"
)

func SetupRouter(appServices *app.AppServices) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	v1 := r.Group("/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("/place", PlaceOrderHandler(appServices))
			orders.POST("/hydrate-orderbook", HydrateOrderbookInRedis(appServices))
		}
	}

	return r
}
