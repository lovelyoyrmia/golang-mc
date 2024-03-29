CREATE TABLE "users" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "uid" varchar UNIQUE NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar,
  "password" varchar NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "secret_key" varchar NOT NULL,
  "is_active" bool DEFAULT true,
  "is_verified" bool DEFAULT false,
  "last_login" timestamp,
  "date_joined" timestamp DEFAULT 'now()'
);

CREATE TABLE "user_otps" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "uid" varchar UNIQUE NOT NULL,
  "otp_enabled" bool DEFAULT false,
  "otp_verified" bool DEFAULT false,
  "otp_secret" varchar,
  "otp_url" varchar
);

CREATE TABLE "sessions" (
  "id" varchar PRIMARY KEY NOT NULL,
  "uid" varchar NOT NULL,
  "refresh_token" varchar,
  "user_agent" varchar,
  "client_ip" varchar,
  "is_blocked" bool DEFAULT false,
  "expired_at" timestamp,
  "created_at" timestamp DEFAULT 'now()'
);

CREATE TABLE "verify_emails" (
  "id" varchar PRIMARY KEY NOT NULL,
  "uid" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_key" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT 'now()',
  "expired_at" timestamp
);

CREATE TABLE "recover_accounts" (
  "id" varchar PRIMARY KEY NOT NULL,
  "uid" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_key" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT 'now()',
  "expired_at" timestamp
);

CREATE INDEX ON "user_otps" ("uid");
