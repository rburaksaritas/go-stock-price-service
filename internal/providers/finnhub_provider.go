package providers

import (
	"encoding/json"
	"fmt"
	"go-stock-price-service/internal/errors"
	"go-stock-price-service/internal/models"
	"go-stock-price-service/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
)

type FinnhubProvider struct {
	apiKey string
}

func NewFinnhubProvider() *FinnhubProvider {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY not set in .env file")
	}
	return &FinnhubProvider{apiKey: apiKey}
}

func (p *FinnhubProvider) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	if stockId == "" {
		return nil, &errors.InvalidInputError{Param: "stockId", Value: stockId}
	}

	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stockId, p.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error making request to Finnhub: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, &errors.ExternalAPIError{Message: "Too Many Requests"}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error response from Finnhub: %s", resp.Status)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error reading response body: %v", err)}
	}

	var rawData models.RawDataFinnhub
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error parsing JSON response: %v", err)}
	}

	readableTimestamp, err := utils.Int64ToReadableTimestamp(rawData.Timestamp, timeZone)
	if err != nil {
		return nil, &errors.InvalidInputError{Param: "timezone", Value: timeZone}
	}

	data := &models.PriceData{
		CurrentPrice:  rawData.CurrentPrice,
		OpenPrice:     rawData.OpenPrice,
		HighPrice:     rawData.HighPrice,
		LowPrice:      rawData.LowPrice,
		PreviousClose: rawData.PreviousClose,
		Timestamp:     *readableTimestamp,
		Provider:      "Finnhub",
	}

	if data.CurrentPrice == 0 && data.OpenPrice == 0 && data.HighPrice == 0 && data.LowPrice == 0 && data.PreviousClose == 0 {
		return nil, &errors.NotFoundError{StockID: stockId}
	}

	return data, nil
}
