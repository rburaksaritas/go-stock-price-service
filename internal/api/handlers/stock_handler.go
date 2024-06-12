package handlers

import (
	"go-stock-price-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	StockService *service.StockService
}

func NewStockHandler(stockService *service.StockService) *StockHandler {
	return &StockHandler{StockService: stockService}
}

func (h *StockHandler) GetStockPrice(c *gin.Context) {
	stockId := c.Param("stockId")
	timeZone := c.DefaultQuery("timezone", "UTC")

	priceData, err := h.StockService.FetchPrice(stockId, timeZone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if priceData.CurrentPrice == 0 && priceData.OpenPrice == 0 && priceData.HighPrice == 0 && priceData.LowPrice == 0 && priceData.PreviousClose == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid stock ID or no data available"})
		return
	}

	c.JSON(http.StatusOK, priceData)
}
