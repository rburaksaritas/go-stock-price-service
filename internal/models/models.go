package models

type PriceData struct {
	CurrentPrice  float64 `json:"current_price"`
	OpenPrice     float64 `json:"open_price"`
	HighPrice     float64 `json:"high_price"`
	LowPrice      float64 `json:"low_price"`
	PreviousClose float64 `json:"previous_close"`
	Timestamp     string  `json:"timestamp"`
}

type RawDataFinnhub struct {
	CurrentPrice  float64 `json:"c"`
	OpenPrice     float64 `json:"o"`
	HighPrice     float64 `json:"h"`
	LowPrice      float64 `json:"l"`
	PreviousClose float64 `json:"pc"`
	Timestamp     int64   `json:"t"`
}