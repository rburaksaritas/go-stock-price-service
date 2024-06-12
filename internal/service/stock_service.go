package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go-stock-price-service/internal/models"
	"go-stock-price-service/pkg/utils"
)

type StockService struct {
	ApiKey string
}

type PriceData struct {
	CurrentPrice  float64 `json:"current_price"`
	OpenPrice     float64 `json:"open_price"`
	HighPrice     float64 `json:"high_price"`
	LowPrice      float64 `json:"low_price"`
	PreviousClose float64 `json:"previous_close"`
	Timestamp     string  `json:"timestamp"`
}

func NewStockService() *StockService {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY not set in .env file")
	}
	return &StockService{ApiKey: apiKey}
}

func (s *StockService) GetStockPrice(c *gin.Context) {
	stockId := c.Param("stockId")
	timeZone := c.DefaultQuery("timezone", "UTC")

	priceData, err := s.fetchPrice(stockId, timeZone)
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

func (s *StockService) fetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stockId, s.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Finnhub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from Finnhub: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var rawData models.RawDataFinnhub
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	readableTimestamp, err := utils.ReadableTimestamp(rawData.Timestamp, timeZone)
	if err != nil {
		return nil, fmt.Errorf("error converting timestamp: %v", err)
	}

	data := &models.PriceData{
		CurrentPrice:  rawData.CurrentPrice,
		OpenPrice:     rawData.OpenPrice,
		HighPrice:     rawData.HighPrice,
		LowPrice:      rawData.LowPrice,
		PreviousClose: rawData.PreviousClose,
		Timestamp:     *readableTimestamp,
	}

	return data, nil
}
