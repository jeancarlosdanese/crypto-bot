-- migrations/0001_create_initial_schema.sql

-- Extensões necessárias
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ================================
-- Tabelas principais
-- ================================

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    whatsapp VARCHAR(20),
    api_key VARCHAR(64) NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    binance_api_key VARCHAR(100),
    binance_api_secret VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX accounts_email_key ON accounts (email);
CREATE UNIQUE INDEX accounts_whatsapp_key ON accounts (whatsapp);

CREATE TABLE account_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    otp_code VARCHAR(8) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    expires_at TIMESTAMP NOT NULL
);
CREATE INDEX account_otps_account_id_idx ON account_otps (account_id);

CREATE TABLE strategies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE bots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    strategy_id UUID NOT NULL REFERENCES strategies(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    interval VARCHAR(10) NOT NULL,
    autonomous BOOLEAN DEFAULT false,
    config_json JSONB DEFAULT '{}'::jsonb,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
CREATE INDEX bots_account_id_idx ON bots (account_id);
CREATE INDEX bots_strategy_id_idx ON bots (strategy_id);

CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id UUID NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    entry_price NUMERIC(18,8) NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
CREATE UNIQUE INDEX positions_bot_id_key ON positions (bot_id);

CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id UUID NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    entry_price NUMERIC(18,8),
    entry_time BIGINT,
    exit_price NUMERIC(18,8),
    exit_time BIGINT,
    duration INT,
    profit NUMERIC(18,8),
    roi_pct NUMERIC(8,4),
    strategy VARCHAR(50),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE decisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id UUID NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    symbol VARCHAR(20),
    interval VARCHAR(10),
    timestamp BIGINT,
    decision VARCHAR(10),
    price NUMERIC(18,8),
    indicators JSONB,
    context JSONB,
    strategy VARCHAR(50),
    created_at TIMESTAMP DEFAULT now()
);
