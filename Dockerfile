FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /stock-price-service ./cmd/server

EXPOSE 8080

CMD [ "/stock-price-service" ]
