-- migrations/0002_create_bots_tables.sql

-- Tabela de bots
CREATE TABLE "public"."bots" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "account_id" uuid NOT NULL,
    "symbol" varchar(20) NOT NULL,
    "interval" varchar(10) NOT NULL,
    "strategy_name" varchar(50) NOT NULL,
    "autonomous" boolean DEFAULT false,
    "active" boolean DEFAULT true,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    CONSTRAINT "bots_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE CASCADE
);
CREATE INDEX bots_account_id_idx ON public.bots USING btree (account_id);

-- Configurações dinâmicas por bot
CREATE TABLE "public"."bot_configs" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "bot_id" uuid NOT NULL,
    "config_json" jsonb NOT NULL,
    "created_at" timestamp DEFAULT now(),
    CONSTRAINT "bot_configs_bot_id_fkey" FOREIGN KEY ("bot_id") REFERENCES "public"."bots"("id") ON DELETE CASCADE
);

-- Posições em aberto
CREATE TABLE "public"."positions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "bot_id" uuid NOT NULL,
    "entry_price" numeric(18,8) NOT NULL,
    "timestamp" bigint NOT NULL,
    CONSTRAINT "positions_bot_id_fkey" FOREIGN KEY ("bot_id") REFERENCES "public"."bots"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX positions_bot_id_key ON public.positions (bot_id);

-- Execuções completas
CREATE TABLE "public"."executions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "bot_id" uuid NOT NULL,
    "entry_price" numeric(18,8),
    "entry_time" bigint,
    "exit_price" numeric(18,8),
    "exit_time" bigint,
    "duration" int,
    "profit" numeric(18,8),
    "roi_pct" numeric(8,4),
    "strategy" varchar(50),
    "created_at" timestamp DEFAULT now(),
    CONSTRAINT "executions_bot_id_fkey" FOREIGN KEY ("bot_id") REFERENCES "public"."bots"("id") ON DELETE CASCADE
);

-- Decisões tomadas
CREATE TABLE "public"."decisions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "bot_id" uuid NOT NULL,
    "symbol" varchar(20),
    "interval" varchar(10),
    "timestamp" bigint,
    "decision" varchar(10),
    "price" numeric(18,8),
    "indicators" jsonb,
    "context" jsonb,
    "strategy" varchar(50),
    "created_at" timestamp DEFAULT now(),
    CONSTRAINT "decisions_bot_id_fkey" FOREIGN KEY ("bot_id") REFERENCES "public"."bots"("id") ON DELETE CASCADE
);