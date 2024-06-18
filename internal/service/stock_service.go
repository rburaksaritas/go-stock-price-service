package service

import (
	"go-stock-price-service/internal/errors"
	"go-stock-price-service/internal/models"
	"go-stock-price-service/internal/providers"
	"go-stock-price-service/pkg/utils"
	"time"
)

type StockService struct {
	providers []providers.Provider
	cache     *utils.RedisClient
}

func NewStockService(providers []providers.Provider, cache *utils.RedisClient) *StockService {
	return &StockService{
		providers: providers,
		cache:     cache,
	}
}

func (s *StockService) FetchPrice(stockId string, timeZone string) (*models.PriceData, error) {
	cacheKey := stockId + "_" + timeZone

	// Check the TTL of the cache key
	ttl, err := s.cache.TTL(cacheKey)
	if err == nil && ttl.Seconds() > 50 {
		var cachedData models.PriceData
		if cacheErr := s.cache.Get(cacheKey, &cachedData); cacheErr == nil {
			cachedData.Provider += " (cached)"
			return &cachedData, nil
		}
	}

	for _, provider := range s.providers {
		data, err := provider.FetchPrice(stockId, timeZone)
		if err == nil {
			// Store the result in the cache
			if cacheErr := s.cache.Set(cacheKey, data, 1*time.Minute); cacheErr != nil {
				return nil, cacheErr
			}
			return data, nil
		}

		if apiErr, ok := err.(*errors.ExternalAPIError); ok && apiErr.Message == "Too Many Requests" {
			continue
		} else {
			return nil, err
		}
	}

	var cachedData models.PriceData
	if cacheErr := s.cache.Get(cacheKey, &cachedData); cacheErr == nil {
		cachedData.Provider += " (cached)"
		return &cachedData, nil
	}

	return nil, &errors.ExternalAPIError{Message: "All providers rate-limited and data is not available in the server-side cache. Please try again later."}
}
