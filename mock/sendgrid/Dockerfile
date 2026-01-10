FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY mock/sendgrid/go.mod mock/sendgrid/go.sum* ./
RUN go mod download

COPY mock/sendgrid/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /sendgrid-mock main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /sendgrid-mock .

EXPOSE 8081

CMD ["./sendgrid-mock"]
