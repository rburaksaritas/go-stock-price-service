# Go Stock Price Service

This is a Go-based stock price service that fetches real-time stock prices from multiple providers and caches the results for improved performance.

## Features

- Fetches stock prices from multiple providers (Finnhub, Alpha Vantage, and Polygon)
- Implements a fallback mechanism when rate limits are reached
- Caches results using Redis for faster subsequent requests
- Supports different time zones for timestamp conversion
- Provides a RESTful API endpoint for retrieving stock prices

## Prerequisites

- Go 1.16 or higher
- Redis server
- API keys for Finnhub, Alpha Vantage, and Polygon

## Installation

1. Clone the repository:
```
git clone https://github.com/rburaksaritas/go-stock-price-service.git
```
2. Navigate to repository:
```
cd go-stock-price-service
```
3. Install dependencies:
```
go mod tidy
```
5. Set up environment variables:
Create a `.env` file in the project root and add the following:
```
FINNHUB_API_KEY=your_finnhub_api_key
ALPHAVANTAGE_API_KEY=your_alphavantage_api_key
POLYGON_API_KEY=your_polygon_api_key
```
## Usage

1. Start the Redis server.

```
brew services start redis
```

2. Run the application:

```
go run cmd/main.go
```

3. The server will start at `http://localhost:8080`.

## API Endpoint

- `GET /prices/:stockId`: Retrieves the current stock price for the given stock ID.
- Query Parameters:
 - `timezone` (optional): Specify the timezone for the timestamp (default: UTC)
- Example: `http://localhost:8080/prices/AAPL?timezone=us`

## Project Structure

- `cmd/main.go`: Entry point of the application
- `api/`: API routing and setup
    - `handlers/`: HTTP request handlers
- `internal/`: Contains the core application code
    - `errors/`: Custom error types
    - `models/`: Data models
    - `providers/`: Stock price data providers
    - `service/`: Business logic for fetching stock prices
- `pkg/`: Reusable packages
    - `utils/`: Utility functions and Redis client

## Error Handling

The service handles various error scenarios, including:
- Invalid input
- Stock not found
- External API errors
- Rate limiting

## Caching

Stock price data is cached in Redis for 1 minute to reduce the number of API calls and improve response times.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
