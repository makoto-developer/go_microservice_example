FROM golang:1.25-alpine AS builder

WORKDIR /app

# manual/authをコピー
COPY manual/auth ./manual/auth

# auth-serviceのファイルをコピー
COPY generated/auth-service/go.mod generated/auth-service/go.sum* ./
RUN go mod download

COPY generated/auth-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /service .

EXPOSE 50051

CMD ["./service"]
