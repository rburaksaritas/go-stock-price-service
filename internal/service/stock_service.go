package service

import (
	"go-stock-price-service/internal/errors"
	"go-stock-price-service/internal/models"
	"go-stock-price-service/internal/providers"
)

type StockService struct {
	finnhubProvider      *providers.FinnhubProvider
	alphaVantageProvider *providers.AlphaVantageProvider
}

func NewStockService(finnhubProvider *providers.FinnhubProvider, alphaVantageProvider *providers.AlphaVantageProvider) *StockService {
	return &StockService{
		finnhubProvider:      finnhubProvider,
		alphaVantageProvider: alphaVantageProvider,
	}
}

func (s *StockService) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	data, err := s.finnhubProvider.FetchPrice(stockId, timeZone)
	if err != nil {
		if apiErr, ok := err.(*errors.ExternalAPIError); ok && apiErr.Message == "Too Many Requests" {
			// Finnhub rate-limited, switch to Alpha Vantage
			data, err = s.alphaVantageProvider.FetchPrice(stockId, timeZone)
			if err != nil {
				if avErr, ok := err.(*errors.ExternalAPIError); ok && avErr.Message == "Too Many Requests" {
					// Alpha Vantage rate-limited, indicate fallback to cached data
					return nil, &errors.ExternalAPIError{Message: "Both providers rate-limited, attempting to provide cached data."}
				}
				return nil, err
			}
		} else {
			// Return other errors as is
			return nil, err
		}
	}
	return data, nil
}
