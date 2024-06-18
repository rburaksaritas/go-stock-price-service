package providers

import "go-stock-price-service/internal/models"

type Provider interface {
	FetchPrice(stockId string, timeZone string) (*models.PriceData, error)
}
