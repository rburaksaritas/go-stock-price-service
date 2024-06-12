#!/bin/bash

# Build the project
go build -o stock-price-service ./cmd/server

# Run the built executable
./stock-price-service
