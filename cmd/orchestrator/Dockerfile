FROM golang:1.23.7 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/orchestrator
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/orchestrator .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/orchestrator .

CMD ["./orchestrator"]