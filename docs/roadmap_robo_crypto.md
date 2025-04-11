# 🧠 Robo de Crypto – Roadmap Atualizado

Este documento organiza e prioriza os próximos passos para evolução do projeto "Robô de Crypto", com foco em robustez, inteligência, modularidade e experiência do usuário.

---

## ✅ Fase Atual – Estrutura Robusta e Testes Locais

### 🔄 Conexão & Execução

- [x] Suporte a múltiplos pares com WebSocket dinâmico
- [x] Reconexão automática e heartbeat
- [x] Controle individual por par (stream isolado por `go routine`)
- [x] Registro de decisões e execuções no MongoDB

### 📊 Estratégia & Logs

- [x] Estratégias inteligentes com MA, RSI, Volume, ATR
- [x] Implementação de trailing stop técnico
- [x] Decisão por slope das EMAs (EMA Fan)
- [x] Logs salvos por operação e por par
- [x] Resumos automáticos após cada trade (`.log`)

---

## 🧱 Fase 2 – Refatoração & Arquitetura Moderna

### ⚙️ Estratégias & Configuração por Bot (por par)

- [x] Estrutura para múltiplos bots rodando simultaneamente
- [x] Cada par com sua `StrategyConfig`
- [ ] Permitir diferentes estratégias por par (`EvaluateCrossover`, `EMA Fan`, `RSI2`, etc.)
- [ ] Configurações dinâmicas por estratégia (parametrizadas)

### 🗃 Banco de Dados & Infra

- [x] Migrar de MongoDB para **PostgreSQL**
- [x] Criar estrutura relacional para:
  - Users
  - Bots
  - Configs
  - Logs
  - Execuções
- [x] Criar scripts de seed e migração

---

## 🖥️ Fase 3 – App Web (Next.js)

### 👤 Multiusuário & Painel de Controle

- [ ] Sistema de contas de usuário
- [ ] Cadastro/login com proteção (JWT ou OAuth)
- [ ] Tela de criação e configuração de Bots:
  - Seleção do par
  - Intervalo
  - Estratégia
  - Parâmetros customizados
  - Autonomia: **manual** ou **automática**
- [ ] Monitoramento em tempo real:
  - Status dos Bots
  - Decisões sobre o gráfico
  - Execuções passadas
  - Logs de performance

---

## 🚀 Fase 4 – Módulo de Execução Real & Risk Management

### 🛠️ Integração com Binance para execução real

- [ ] Criar `BinanceTradeService` com:
  - Envio de ordens reais
  - Consulta de saldo e posições
  - Cancelamento de ordens

### 🛡️ Gestão de Risco

- [ ] Stop-loss e take-profit customizáveis
- [ ] Definição de tamanho de posição por risco ou capital
- [ ] Trava de emergência
- [ ] Modo Paper Trading (simulado)

---

## 🔔 Fase 5 – Alertas, IA e Backtesting

### 🔔 Alertas e Automação

- [ ] Alertas via Telegram, e-mail ou webhook
- [ ] Painel com notificações de erro/sinal

### 🧠 IA, Otimização e Aprendizado

- [ ] Log completo com estratégia usada (nome, parâmetros, resultado)
- [ ] Painel de comparação entre estratégias
- [ ] Módulo de backtesting com simulação completa

---

## 📌 Prioridades imediatas sugeridas

1. ✅ Manter testes com estratégia atual rodando local (BTC, ETH, SOL)
2. ⚙️ Criar `StrategyConfig` por par com injeção no `main.go`
3. 🛠️ Começar migração para PostgreSQL
4. 🖥️ Planejar modelo de dados para multiusuários + API Keys
5. 🚀 Planejar arquitetura do App Web em Next.js (com ou sem painel real-time)
