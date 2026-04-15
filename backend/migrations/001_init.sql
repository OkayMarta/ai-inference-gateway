-- Базова схема для Lab 3: PostgreSQL persistence.

BEGIN;

-- Користувачі системи та їхній токен-баланс.
CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    username VARCHAR NOT NULL,
    token_balance NUMERIC(12,2) NOT NULL CHECK (token_balance >= 0)
);

-- Доступні AI-моделі та їхня вартість у токенах.
CREATE TABLE ai_models (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    token_cost NUMERIC(12,2) NOT NULL CHECK (token_cost >= 0)
);

-- Завдання на обробку prompt payload із фіксацією стану виконання.
CREATE TABLE prompt_tasks (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    model_id VARCHAR NOT NULL REFERENCES ai_models(id) ON DELETE RESTRICT,
    payload TEXT NOT NULL,
    status VARCHAR NOT NULL,
    result TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_prompt_tasks_status
        CHECK (status IN ('Queued', 'Processing', 'Completed', 'Failed', 'Cancelled'))
);

-- Фінансові операції користувача: списання та повернення.
CREATE TABLE transactions (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    task_id VARCHAR REFERENCES prompt_tasks(id) ON DELETE SET NULL,
    amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
    type VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_transactions_type
        CHECK (type IN ('charge', 'refund'))
);

-- Воркери, що виконують задачі та можуть масштабуватись окремо в майбутньому.
CREATE TABLE worker_nodes (
    id VARCHAR PRIMARY KEY,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_worker_nodes_status
        CHECK (status IN ('Idle', 'Busy'))
);

-- Зв'язок many-to-many між воркерами та підтримуваними моделями.
CREATE TABLE worker_supported_models (
    worker_id VARCHAR NOT NULL REFERENCES worker_nodes(id) ON DELETE CASCADE,
    model_id VARCHAR NOT NULL REFERENCES ai_models(id) ON DELETE CASCADE,
    PRIMARY KEY (worker_id, model_id)
);

-- Індекси для типових сценаріїв читання: черги, фільтрації та історії.
CREATE INDEX idx_prompt_tasks_user_id ON prompt_tasks(user_id);
CREATE INDEX idx_prompt_tasks_status ON prompt_tasks(status);
CREATE INDEX idx_prompt_tasks_created_at ON prompt_tasks(created_at);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);

CREATE INDEX idx_worker_supported_models_model_id ON worker_supported_models(model_id);

COMMIT;
