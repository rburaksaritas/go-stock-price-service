package handlers

import (
	"go-stock-price-service/internal/errors"
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
		switch e := err.(type) {
		case *errors.InvalidInputError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		case *errors.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.ExternalAPIError:
			c.JSON(http.StatusBadGateway, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error has occured."})
		}
		return
	}

	c.JSON(http.StatusOK, priceData)
}
