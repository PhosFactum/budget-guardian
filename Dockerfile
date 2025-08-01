# Build Stage
FROM golang:1.24.4-alpine AS builder

# Installing git for Go-alpine
RUN apk add --no-cache git ca-certificates

# Configuring workdir and dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o budget-guardian ./cmd/bot/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/budget-guardian .
COPY --from=builder /app/.env .

# Exposing port and running app
EXPOSE 8080
CMD ["./budget-guardian"]