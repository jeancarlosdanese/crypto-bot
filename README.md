# Crypto Trading Bot

Este projeto é um robô de trading automatizado desenvolvido em Go, utilizando Clean Architecture, DDD, MongoDB e Docker. Ele se conecta à Binance e executa operações com base em uma estratégia própria, com foco em preservação e multiplicação de capital em dólar.

## 📌 Estratégia

A lógica completa da estratégia pode ser consultada no arquivo [STRATEGY.md](./STRATEGY.md). Ela inclui:

- Classificação dos ativos da carteira (reservas, especulativos, alto risco)
- Decisões baseadas em análise de mercado e volume
- Gestão de risco e controle de lucros/perdas

## ⚙️ Tecnologias

- Go
- MongoDB
- Docker e Docker Compose
- Arquitetura Limpa (Clean Architecture)
- DDD (Domain-Driven Design)

## 📁 Estrutura do Projeto

```
.
├── cmd/api                  # Entry point da aplicação (main.go)
├── internal
│   ├── app                 # Casos de uso e serviços
│   ├── domain              # Entidades e contratos de domínio
│   ├── infra               # Infraestrutura (MongoDB, config, etc.)
│   └── interfaces          # Interfaces externas (handlers, controllers)
├── test                    # Testes unitários e mocks
├── configs                 # Arquivos de configuração
├── scripts                 # Scripts utilitários
├── STRATEGY.md             # Estratégia de trading descrita em detalhes
├── docker-compose.yaml     # Ambiente de desenvolvimento com MongoDB
├── .env                    # Variáveis de ambiente
├── Makefile                # Comandos úteis (build, test, cover etc.)
```

## ▶️ Executando Localmente

1. Clone o repositório
2. Copie o `.env.example` para `.env` e configure com suas chaves da Binance
3. Suba os serviços com Docker:

```bash
docker-compose up -d
```

4. Execute o projeto:

```bash
go run ./cmd/api/main.go
```

## ✅ Testes

```bash
make test      # Executa testes
make cover     # Gera relatório de cobertura
```

## 📌 Licença

Este projeto é open source e está licenciado sob os termos da [MIT License](LICENSE).
