package service

import (
	"encoding/json"
	"fmt"
	"go-stock-price-service/internal/models"
	"go-stock-price-service/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
)

type StockService struct {
	ApiKey string
}

func NewStockService() *StockService {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY not set in .env file")
	}
	return &StockService{ApiKey: apiKey}
}

func (s *StockService) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
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
