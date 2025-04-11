# 📦 Checkpoint: Projeto Robô de Crypto

> Última atualização: 2025-04-11 11:11:45

---

## ✅ Status Geral

- Projeto em execução local com **Docker Compose**
- Banco de dados: **PostgreSQL**
- Frontend Web: ainda não iniciado
- Estratégia atual: `EvaluateCrossover`
- Execução com múltiplos bots (BTC, ETH, SOL, BNB, XRP, FDUSD)
- Configuração dinâmica por bot no banco (`bot_configs`)

---

## ⚙️ Infraestrutura

### Docker Compose

- Container `postgres` com volume persistente `pg_data`
- Container `crypto-bot` compilado com `golang:1.23`, imagem final baseada em `alpine`
- Conexão verificada via `psql` com `wait-for-postgres.sh`
- Migrations aplicadas automaticamente (`/migrations`)
- `.env` com as variáveis necessárias:

```
POSTGRES_USER=hyberica
POSTGRES_PASSWORD=password
POSTGRES_HOST=postgres
POSTGRES_DATABASE=crypto-bot
POSTGRES_PORT=5432

BINANCE_API_KEY=key
BINANCE_API_SECRET=secret
```

---

## 🧠 Estratégia Atual

### Estratégia `EvaluateCrossover`:

- Baseada em:
  - MA(9), MA(26)
  - RSI (limiar 70)
  - ATR para trailing stop
  - EMA(5) para saída técnica
- Entrada: MA9 > MA26, RSI < 70, preço > MA9
- Saída: ATR Stop, RSI reversão, preço < EMA5, sinal de venda

---

## 📊 Banco de Dados

### Tabelas principais:

- `accounts` (usuários)
- `bots` (por par e estratégia)
- `bot_configs` (configurações por bot em JSONB)
- `positions` (posição em aberto)
- `executions` (execuções finalizadas)
- `decisions` (decisões registradas com contexto e indicadores)

---

## 🔧 Scripts

### `wait-for-postgres.sh`
Espera o PostgreSQL estar pronto antes de iniciar o bot.

---

## ✅ Concluído recentemente

- Migração completa de MongoDB para PostgreSQL
- Docker funcional com builds em `alpine`
- Estratégias com logs detalhados e performance por minuto
- Logs `.log` por minuto com estatísticas
- Suporte a múltiplos pares rodando em paralelo

---

## ⏭️ Próximos passos sugeridos

1. Permitir múltiplas estratégias (`EMA Fan`, `RSI`, etc.)
2. Interface Web com Next.js (painel e controle de bots)
3. Integração com Binance real
4. Painel com IA, comparações e alertas

---

_Ficheiro gerado automaticamente para checkpoint do projeto Robô de Crypto._
