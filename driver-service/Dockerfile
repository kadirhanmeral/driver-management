FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o driver-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/driver-service .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

RUN chmod +x driver-service

EXPOSE 8081

CMD ["./driver-service"]
