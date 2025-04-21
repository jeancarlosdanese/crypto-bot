-- migrations/9999_init_example.sql

-- ================================
-- Dados iniciais
-- ================================

-- Estratégias disponíveis
INSERT INTO strategies (name, description) VALUES
('CROSSOVER', 'Cruzamento de Médias Móveis com RSI'),
('EMA_FAN', 'Alinhamento de múltiplas EMAs com volume'),
('RSI2', 'Estratégia baseada no RSI2 e reversão técnica');

-- Conta administrativa
INSERT INTO accounts (
    id, name, email, whatsapp, api_key, binance_api_key, binance_api_secret
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'João da Silva',
    'email@domain.com',
    '9999999999',
    encode(gen_random_bytes(32), 'hex'),
    'SUA_API_KEY',
    'SUA_API_SECRET'
);

-- ================================
-- Inserção de bots de exemplo
-- ================================
DO $$
DECLARE
    acc_id UUID := (SELECT id FROM accounts WHERE email = 'email@domain.com');
    crossover_id UUID := (SELECT id FROM strategies WHERE name = 'CROSSOVER');
    symbol TEXT;
    bot_id UUID;
BEGIN
    FOR symbol IN SELECT unnest(ARRAY['BTC/USDT', 'BNB/USDT', 'XRP/USDT', 'ETH/USDT', 'SOL/USDT', 'FDUSD/USDT'])
    LOOP
        INSERT INTO bots (
            id, account_id, symbol, interval,
            strategy_id, autonomous, active, config_json
        )
        VALUES (
            gen_random_uuid(),
            acc_id,
            symbol,
            '1m',
            crossover_id,
            true,
            true,
            '{
                "ema_periods": [9, 26],
                "macd": { "short": 12, "long": 26, "signal": 9 },
                "rsi_period": 14,
                "rsi_buy": 10,
                "rsi_sell": 90,
                "bollinger": { "period": 20 },
                "atr_period": 14,
                "volatility_window": 14
            }'::jsonb
        )
        RETURNING id INTO bot_id;
    END LOOP;
END $$;
