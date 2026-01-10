FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY mock/carriers/go.mod mock/carriers/go.sum* ./
RUN go mod download
COPY mock/carriers/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /carriers-mock main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /carriers-mock .
EXPOSE 8083
CMD ["./carriers-mock"]
