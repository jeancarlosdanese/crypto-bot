-- migrations/0001_create_acconts_table.sql

-- Table Definition
CREATE TABLE "public"."accounts" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" varchar(100) NOT NULL,
    "email" varchar(150) NOT NULL,
    "whatsapp" varchar(20),
    "is_admin" boolean DEFAULT false,
    "api_key" varchar(64),
    "binance_api_key" varchar(100),
    "binance_api_secret" varchar(100),
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Indices
CREATE UNIQUE INDEX accounts_email_key ON public.accounts USING btree (email);
CREATE UNIQUE INDEX accounts_whatsapp_key ON public.accounts USING btree (whatsapp);


CREATE TABLE "public"."account_otps" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "account_id" uuid,
    "otp_code" varchar(8) NOT NULL,
    "created_at" timestamp DEFAULT now(),
    "expires_at" timestamp NOT NULL,
    CONSTRAINT "account_otps_account_id_fkey"
        FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);
CREATE INDEX account_otps_account_id_idx ON public.account_otps USING btree (account_id);
