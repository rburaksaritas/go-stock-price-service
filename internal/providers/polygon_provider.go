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
	"strconv"
)

type PolygonProvider struct {
	apiKey string
}

func NewPolygonProvider() *PolygonProvider {
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		log.Fatal("POLYGON_API_KEY not set in .env file")
	}
	return &PolygonProvider{apiKey: apiKey}
}

func (p *PolygonProvider) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?adjusted=true&apiKey=%s", stockId, p.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error making request to Polygon: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, &errors.ExternalAPIError{Message: "Too Many Requests"}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error response from Polygon: %s", resp.Status)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error reading response body: %v", err)}
	}

	var response struct {
		Results []struct {
			Close     float64         `json:"c"`
			High      float64         `json:"h"`
			Low       float64         `json:"l"`
			Open      float64         `json:"o"`
			Volume    json.Number     `json:"v"`
			Timestamp json.RawMessage `json:"t"` // Use json.RawMessage to handle different types
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error parsing JSON response: %v", err)}
	}

	if len(response.Results) == 0 {
		return nil, &errors.NotFoundError{StockID: stockId}
	}

	latestData := response.Results[0]

	var timestamp int64
	if err := json.Unmarshal(latestData.Timestamp, &timestamp); err != nil {
		var timestampStr string
		if err := json.Unmarshal(latestData.Timestamp, &timestampStr); err != nil {
			return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error converting timestamp: %v", err)}
		}
		timestamp, err = strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return nil, &errors.ExternalAPIError{Message: fmt.Sprintf("Error converting timestamp: %v", err)}
		}
	}

	readableTimestamp, err := utils.Int64ToReadableTimestamp(timestamp, timeZone)
	if err != nil {
		return nil, &errors.InvalidInputError{Param: "timezone", Value: timeZone}
	}

	data := &models.PriceData{
		CurrentPrice:  latestData.Close,
		OpenPrice:     latestData.Open,
		HighPrice:     latestData.High,
		LowPrice:      latestData.Low,
		PreviousClose: latestData.Close,
		Timestamp:     *readableTimestamp,
		Provider:      "Polygon",
	}

	return data, nil
}
