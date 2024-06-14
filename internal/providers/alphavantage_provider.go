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
	"time"
)

type AlphaVantageProvider struct {
	apiKey string
}

func NewAlphaVantageProvider() *AlphaVantageProvider {
	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("ALPHAVANTAGE_API_KEY not set in .env file")
	}
	return &AlphaVantageProvider{apiKey: apiKey}
}

func (p *AlphaVantageProvider) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", stockId, p.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error making request to Alpha Vantage: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, &errors.ExternalAPIError{Message: "Too Many Requests"}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error response from Alpha Vantage: %s", resp.Status)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error reading response body: %v", err)}
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error parsing JSON response: %v", err)}
	}

	if _, ok := responseMap["Note"]; ok && len(responseMap) == 1 {
		return nil, &errors.ExternalAPIError{Message: "Too Many Requests"}
	}

	if _, ok := responseMap["Information"]; ok && len(responseMap) == 1 {
		return nil, &errors.ExternalAPIError{Message: "Too Many Requests"}
	}

	var rawData models.RawDataAlphaVantage
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error parsing JSON response: %v", err)}
	}

	readableTimestamp, err := utils.TimeToReadableTimestamp(time.Now().UTC(), timeZone)
	if err != nil {
		return nil, &errors.InvalidInputError{Param: "timezone", Value: timeZone}
	}

	data := &models.PriceData{
		CurrentPrice:  rawData.GlobalQuote.Price,
		OpenPrice:     rawData.GlobalQuote.Open,
		HighPrice:     rawData.GlobalQuote.High,
		LowPrice:      rawData.GlobalQuote.Low,
		PreviousClose: rawData.GlobalQuote.PreviousClose,
		Timestamp:     *readableTimestamp,
		Provider:      "Alpha Vantage",
	}

	if data.CurrentPrice == 0 && data.OpenPrice == 0 && data.HighPrice == 0 && data.LowPrice == 0 && data.PreviousClose == 0 {
		return nil, &errors.NotFoundError{StockID: stockId}
	}

	return data, nil
}
