package main

import (
	"go-stock-price-service/internal/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := api.SetupRoutes()

	log.Println("Server started at :8080")
	log.Fatal(router.Run(":8080"))
}
