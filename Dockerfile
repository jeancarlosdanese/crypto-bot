# Etapa 1: build da aplicação
FROM golang:1.22 AS builder

WORKDIR /app

# Copia os arquivos
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Etapa 2: imagem final enxuta
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/app .

CMD ["/app/app"]
