package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type StockService struct {
	ApiKey string
}

type PriceData struct {
	CurrentPrice  float64 `json:"c"`
	OpenPrice     float64 `json:"o"`
	HighPrice     float64 `json:"h"`
	LowPrice      float64 `json:"l"`
	PreviousClose float64 `json:"pc"`
	Timestamp     int64   `json:"t"`
}

func NewStockService() *StockService {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY not set in .env file")
	}
	return &StockService{ApiKey: apiKey}
}

func (s *StockService) GetStockPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockId := vars["stockId"]
	priceData, err := s.fetchPrice(stockId)
	if err != nil {
		http.Error(w, "Failed to fetch stock price", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(priceData)
}

func (s *StockService) fetchPrice(stockId string) (*PriceData, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stockId, s.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Finnhub: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var data PriceData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &data, nil
}
