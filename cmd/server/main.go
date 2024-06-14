package main

import (
	"go-stock-price-service/internal/api"
	"go-stock-price-service/internal/providers"
	"go-stock-price-service/internal/service"
	"go-stock-price-service/pkg/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cache := utils.NewRedisClient("localhost:6379", "", 0)

	finnhubProvider := providers.NewFinnhubProvider()
	alphaVantageProvider := providers.NewAlphaVantageProvider()
	polygonProvider := providers.NewPolygonProvider()
	stockService := service.NewStockService(finnhubProvider, alphaVantageProvider, polygonProvider, cache)

	router := api.SetupRoutes(stockService)

	log.Println("Server started at :8080")
	log.Fatal(router.Run(":8080"))
}
