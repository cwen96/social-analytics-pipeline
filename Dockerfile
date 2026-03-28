# --- Build stage ---
FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/consumer ./cmd/consumer
RUN go build -o /bin/producer ./cmd/producer
RUN go build -o /bin/api ./cmd/api

# --- Producer image ---
FROM alpine:3.19 AS producer
COPY --from=builder /bin/producer /usr/local/bin/producer
CMD ["producer"]

# --- Consumer image ---
FROM alpine:3.19 AS consumer
COPY --from=builder /bin/consumer /usr/local/bin/consumer
EXPOSE 8080
CMD ["consumer"]

# --- API image ---
FROM alpine:3.19 AS api
COPY --from=builder /bin/api /usr/local/bin/api
EXPOSE 8080
CMD ["api"]
