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

func NewStockService() *StockService {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	log.Println("API Key:", os.Getenv("FINNHUB_API_KEY"))
	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY not set in .env file")
	}
	return &StockService{ApiKey: apiKey}
}

func (s *StockService) GetStockPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockId := vars["stockId"]
	price, err := s.fetchPrice(stockId)
	if err != nil {
		http.Error(w, "Failed to fetch stock price", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"price": price})
}

func (s *StockService) fetchPrice(stockId string) (float64, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stockId, s.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making request to Finnhub: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data struct {
		CurrentPrice float64 `json:"c"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return data.CurrentPrice, nil
}
