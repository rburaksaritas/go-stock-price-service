package main

import (
	"go-stock-price-service/internal/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := mux.NewRouter()
	api.SetupRoutes(router)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
