-- =============================================================
-- Atur Dana — Database Schema
-- Generated from GORM models (internal/models/models.go)
-- PostgreSQL
-- =============================================================

-- gorm.Model provides: id, created_at, updated_at, deleted_at (soft delete)

-- -------------------------------------------------------------
-- users
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,

    username      TEXT         NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    email         TEXT         NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- -------------------------------------------------------------
-- categories
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS categories (
    id         BIGSERIAL    PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    user_id    BIGINT       REFERENCES users (id),
    name       TEXT         NOT NULL UNIQUE,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);
CREATE INDEX IF NOT EXISTS idx_categories_user_id    ON categories (user_id);

-- -------------------------------------------------------------
-- transactions
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS transactions (
    id          BIGSERIAL    PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    user_id     BIGINT       NOT NULL REFERENCES users (id),
    type        TEXT         NOT NULL,
    amount      FLOAT8       NOT NULL,
    description TEXT,
    category_id BIGINT       NOT NULL REFERENCES categories (id),
    date        TIMESTAMP    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions (deleted_at);
-- Composite index for filtered listing (user_id, date, category_id, type)
CREATE INDEX IF NOT EXISTS idx_trx_filter ON transactions (user_id, date, category_id, type);

-- -------------------------------------------------------------
-- budgets
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS budgets (
    id          BIGSERIAL    PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    user_id     BIGINT       NOT NULL REFERENCES users (id),
    category_id BIGINT       NOT NULL REFERENCES categories (id),
    amount      FLOAT8       NOT NULL,
    start_date  TIMESTAMP    NOT NULL,
    end_date    TIMESTAMP    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_budgets_deleted_at ON budgets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_budgets_user_id    ON budgets (user_id);
