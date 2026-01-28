FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY mock/elasticsearch/go.mod mock/elasticsearch/go.sum* ./
RUN go mod download
COPY mock/elasticsearch/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /elasticsearch-mock main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /elasticsearch-mock .
EXPOSE 9200
CMD ["./elasticsearch-mock"]
