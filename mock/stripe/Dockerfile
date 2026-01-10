FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY mock/stripe/go.mod mock/stripe/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY mock/stripe/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /stripe-mock main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /stripe-mock .

EXPOSE 8080

CMD ["./stripe-mock"]
