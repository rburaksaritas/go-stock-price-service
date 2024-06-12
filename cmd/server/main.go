package main

import (
	"go-stock-price-service/internal/api"
	"go-stock-price-service/internal/service"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	stockService := service.NewStockService()
	router := api.SetupRoutes(stockService)

	log.Println("Server started at :8080")
	log.Fatal(router.Run(":8080"))
}
