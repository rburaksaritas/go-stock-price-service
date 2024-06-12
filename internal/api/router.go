package api

import (
	"go-stock-price-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	stockService := service.NewStockService()
	router.GET("/prices/:stockId", stockService.GetStockPrice)
	return router
}
