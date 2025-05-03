FROM golang:1.24.2 AS builder

RUN apt-get update && \
    apt-get install -y gcc libc6-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

WORKDIR /app/cmd/orchestrator
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/orchestrator .

FROM debian:stable-slim

WORKDIR /app

RUN mkdir -p /app/data && \
    chown -R 1000:1000 /app/data

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/orchestrator .

USER 1000
CMD ["./orchestrator"]
