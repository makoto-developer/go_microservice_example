FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod from generated code
COPY generated/auth-service/go.mod generated/auth-service/go.sum* ./
RUN go mod download

# Copy all generated code
COPY generated/auth-service/ ./

# Copy manual implementations if any
COPY manual/auth-service* ./manual/ 2>/dev/null || true

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /auth-service .

EXPOSE 50051

CMD ["./auth-service"]
