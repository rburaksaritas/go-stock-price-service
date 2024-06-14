package service

import (
	"go-stock-price-service/internal/errors"
	"go-stock-price-service/internal/models"
	"go-stock-price-service/internal/providers"
	"go-stock-price-service/pkg/utils"
	"time"
)

type StockService struct {
	finnhubProvider      *providers.FinnhubProvider
	alphaVantageProvider *providers.AlphaVantageProvider
	polygonProvider      *providers.PolygonProvider
	cache                *utils.RedisClient
}

func NewStockService(finnhubProvider *providers.FinnhubProvider, alphaVantageProvider *providers.AlphaVantageProvider, polygonProvider *providers.PolygonProvider, cache *utils.RedisClient) *StockService {
	return &StockService{
		finnhubProvider:      finnhubProvider,
		alphaVantageProvider: alphaVantageProvider,
		polygonProvider:      polygonProvider,
		cache:                cache,
	}
}

func (s *StockService) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	cacheKey := stockId + "_" + timeZone

	data, err := s.finnhubProvider.FetchPrice(stockId, timeZone)
	if err != nil {
		if apiErr, ok := err.(*errors.ExternalAPIError); ok && apiErr.Message == "Too Many Requests" {
			data, err = s.polygonProvider.FetchPrice(stockId, timeZone)
			if err != nil {
				if polyErr, ok := err.(*errors.ExternalAPIError); ok && polyErr.Message == "Too Many Requests" {
					data, err = s.alphaVantageProvider.FetchPrice(stockId, timeZone)
					if err != nil {
						if avErr, ok := err.(*errors.ExternalAPIError); ok && avErr.Message == "Too Many Requests" {
							var cachedData models.PriceData
							if cacheErr := s.cache.Get(cacheKey, &cachedData); cacheErr == nil {
								cachedData.Provider += " (cached)"
								return &cachedData, nil
							}
							return nil, &errors.ExternalAPIError{Message: "All providers rate-limited and data is not available in the server-side cache. Please try again later."}
						}
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	// Store the result in the cache
	if err := s.cache.Set(cacheKey, data, 1*time.Minute); err != nil {
		return nil, err
	}
	return data, nil
}
