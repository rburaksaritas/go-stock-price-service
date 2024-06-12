package api

import (
	"go-stock-price-service/internal/api/handlers"
	"go-stock-price-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(stockService *service.StockService) *gin.Engine {
	router := gin.Default()
	stockHandler := handlers.NewStockHandler(stockService)
	router.GET("/prices/:stockId", stockHandler.GetStockPrice)
	return router
}
