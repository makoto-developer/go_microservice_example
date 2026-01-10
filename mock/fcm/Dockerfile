FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY mock/fcm/go.mod mock/fcm/go.sum* ./
RUN go mod download
COPY mock/fcm/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /fcm-mock main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /fcm-mock .
EXPOSE 8082
CMD ["./fcm-mock"]
