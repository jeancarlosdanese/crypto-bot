# 🧠 Robô de Crypto – Roadmap Atualizado (2025-04-12)

Este documento organiza e prioriza os próximos passos para evolução do projeto "Robô de Crypto", com foco em robustez, inteligência, modularidade e experiência do usuário.

---

## ✅ Fase Atual – Estrutura Robusta e Testes Locais

### 🔄 Conexão & Execução
- [x] Suporte a múltiplos pares com WebSocket dinâmico
- [x] Reconexão automática e heartbeat
- [x] Controle individual por par (stream isolado por `goroutine`)
- [x] Registro de decisões e execuções no PostgreSQL

### 📊 Estratégia & Logs
- [x] Estratégias com MA, RSI, Volume, ATR
- [x] Trailing stop técnico por ATR ou EMA
- [x] Decisão por slope das EMAs (EMA Fan)
- [x] Logs salvos por operação e por par
- [x] Resumos automáticos após cada trade (`.log`)

---

## 🧱 Fase 2 – Refatoração & Arquitetura Moderna

### ⚙️ Estratégias & Configuração por Bot (par)
- [x] Estrutura para múltiplos bots simultâneos
- [x] Configurações por par (via `StrategyConfig`)
- [x] Configurações dinâmicas por bot via banco (JSON)
- [ ] Suporte a múltiplas estratégias (`EvaluateCrossover`, `EMA Fan`, etc.)
- [ ] Injeção da estratégia no `main.go` conforme configuração

### 🗃 Banco de Dados & Infra
- [x] Migração de MongoDB para PostgreSQL
- [x] Tabelas: `accounts`, `bots`, `bot_configs`, `positions`, `executions`, `decisions`
- [x] Scripts SQL de migration e seed
- [x] Docker Compose com Postgres + build Go
- [x] Script `entrypoint.sh` com `wait-for-postgres.sh`

---

## 🖥️ Fase 3 – App Web (Next.js)

### 👤 Multiusuário & Painel de Controle
- [x] Gráficos em tempo real
- [x] Histórico de candles via REST
- [x] Tooltips, médias móveis e decisões visuais
- [x] Dark/Light com detecção dinâmica
- [ ] Autenticação (JWT ou OAuth)
- [ ] Painel de bots por conta
- [ ] Configuração por bot:
  - Par e intervalo
  - Estratégia e parâmetros
  - Autonomia: manual ou automático
- [ ] Histórico de execuções e decisões

---

## 🚀 Fase 4 – Execução Real & Risk Management

### 🛠️ Integração Binance
- [ ] Serviço de ordens reais (BinanceTradeService)
- [ ] Consulta de saldo/posição
- [ ] Modo paper trading
- [ ] Cancelamento automático

### 🛡️ Gestão de Risco
- [ ] Stop-loss/take-profit customizáveis
- [ ] Tamanho de posição baseado em risco
- [ ] Travas e limites

---

## 🔔 Fase 5 – Alertas, IA e Backtesting

### 🔔 Alertas
- [ ] Envio por e-mail, Telegram ou webhook
- [ ] Painel de erros/sinais

### 🧠 IA e Análise
- [x] Log detalhado com estratégia usada (nome, parâmetros, ROI)
- [ ] Benchmark entre estratégias
- [ ] Módulo de backtesting

---

## 📌 Prioridades Imediatas (April/2025)
1. ✅ Finalizar migração para PostgreSQL com bots funcionando
2. ✅ Executar múltiplos bots com configs separadas
3. 🔄 Implementar múltiplas estratégias por bot
4. 🧱 Modularizar criação de bots por estratégia (DI)
5. 🖥️ Iniciar planejamento do app web em Next.js