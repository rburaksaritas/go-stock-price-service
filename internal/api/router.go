package api

import (
	"go-stock-price-service/internal/service"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	stockService := service.NewStockService()
	r.HandleFunc("/prices/{stockId}", stockService.GetStockPrice).Methods("GET")
}
