FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/agent
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/agent .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/agent .
COPY .env .env

USER 1000
CMD ["./agent"]