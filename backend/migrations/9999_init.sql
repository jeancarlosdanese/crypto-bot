-- migrations/init.sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO "public"."accounts" (
    "id", "name", "email", "whatsapp", "api_key", "binance_api_key", "binance_api_secret"
)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Jean Carlos Danese',
    'jean@danese.com.br',
    '5549999669869',
    encode(gen_random_bytes(32), 'hex'), 
    'HB0TGPBzFePBdm7HH40fuYrz54cRZFLYFqkhT0o5fBfB2gFrdzLndsa2oHioE3k6', 
    'sYPkAfh2Trw6zbPzHwUfmvvK2f73qq9bpuZimBoSgT1zq3lcIsg6C8NGaj2hzq7x'
);

-- Substitua 'jean@danese.com.br' pelo e-mail da conta desejada se necessário
DO $$
DECLARE
    acc_id UUID := (SELECT id FROM accounts WHERE email = 'jean@danese.com.br');
    symbol TEXT;
    bot_id UUID;
BEGIN
    -- Insere os bots
    FOR symbol IN 
        SELECT unnest(ARRAY['btcusdt', 'bnbusdt', 'xrpusdt', 'ethusdt', 'solusdt', 'fdusdusdt']) AS symbol
    LOOP
        INSERT INTO bots (id, account_id, symbol, interval, strategy_name, autonomous, active)
        VALUES (
            gen_random_uuid(),
            acc_id,
            symbol,
            '1m',
            'EvaluateCrossover',
            true,
            true
        )
        RETURNING id INTO bot_id;

        INSERT INTO bot_configs (id, bot_id, config_json)
        VALUES (
            gen_random_uuid(),
            bot_id,
            '{"atr_min": 0.5, "ma_long": 26, "ma_short": 9, "rsi_threshold": 70, "volatility_min": 0.1}'::jsonb
        );
    END LOOP;
END $$;

