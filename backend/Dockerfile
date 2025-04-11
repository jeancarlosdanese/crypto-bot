# Etapa 1: Build da aplicação
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-bot ./cmd/api

# Etapa 2: Imagem enxuta, mas com shell
FROM alpine:3.19

WORKDIR /app

# Instala bash e outras ferramentas úteis (opcional, mas prático)
RUN apk add --no-cache bash

COPY --from=builder /app/crypto-bot .

CMD ["/app/crypto-bot"]
