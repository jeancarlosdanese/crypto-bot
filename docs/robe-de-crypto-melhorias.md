# 🧠 Robe de Crypto - Melhorias e Expansões Futuras

Este documento lista sugestões organizadas para evoluir o projeto "Robe de Crypto", com foco em robustez, inteligência e integração.

---

## 🔄 1. Aprimorar o WebSocket e Subscrições

- [x] Implementar reconexão automática após 24h ou desconexões inesperadas.
- [x] Suporte a múltiplos pares (ex: `btcusdt`, `ethusdt`) com subscrição dinâmica.
- [x] Gerenciar ping/pong de forma robusta para evitar desconexões.
- [x] Logging e métricas sobre estabilidade da conexão.

---

## 📊 2. Persistência Inteligente dos Dados

- [x] Armazenar candles finalizados no MongoDB.
- [x] Criar histórico de sinais e decisões para backtesting.
- [x] Estrutura para armazenar logs de execução e performance da estratégia.
- [x] Estado persistente, rastreamento completo de decisões, e execuções à prova de reinício.
- [x] Indexação por tempo e símbolo para consultas eficientes.

---

## ⚙️ 3. Módulo de Execução de Ordens (Binance API)

- [ ] Criar `BinanceTradeService` com suporte a:
  - Enviar ordens reais (compra/venda);
  - Cancelar ordens pendentes;
  - Consultar saldo e posição atual.
- [ ] Implementar modo "simulado" (paper trading) e "real".
- [ ] Validar execução antes de enviar ordens (ex: volatilidade mínima).

---

## 📐 4. Melhorias Estratégicas

- [ ] Detecção automática de topos e fundos.
- [ ] Estratégias de reversão com RSI e ATR.
- [ ] Thresholds adaptativos com base na volatilidade.
- [ ] Implementar filtros baseados em volume ou liquidez.

---

## 🛡️ 5. Gestão de Risco

- [ ] Stop-loss e take-profit configuráveis.
- [ ] Reavaliação dinâmica da posição com base em novos candles.
- [ ] Cálculo de tamanho de posição com base em capital e risco por operação.
- [ ] Trava de emergência em caso de quedas acentuadas ou bugs.

---

## 🚨 6. Alertas e Notificações

- [ ] Enviar alertas via Telegram, Discord ou e-mail:
  - Novos sinais (BUY/SELL)
  - Ordens executadas
  - Erros e reconexões
- [ ] Painel ou log com histórico de alertas.

---

## 📅 7. Roadmap Futuro

- [ ] Criação de um painel web com métricas em tempo real.
- [ ] Modo backtesting completo com simulação de ordens.
- [ ] Estratégias baseadas em aprendizado de máquina (ML).
- [ ] Suporte a múltiplos pares e múltiplas estratégias simultâneas.

## 🧠 Cenário futuro: múltiplas estratégias configuráveis por usuário

Antecipando:

- Que o robô terá diferentes estratégias de decisão (ex: EvaluateCrossover, EvaluateRSI, EvaluateVolatility)
- Que o usuário poderá escolher qual estratégia aplicar
- Que você quer salvar no log qual estratégia foi usada naquela decisão
- E que o log seja útil para comparar estratégias no futuro (benchmark, IA, etc.)
