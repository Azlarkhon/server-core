# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o producer ./examples/main.go

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/producer /app/producer
EXPOSE 8080
CMD ["/app/producer"]