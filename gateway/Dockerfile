FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

RUN chmod +x gateway

EXPOSE 8080

CMD ["./gateway"]
