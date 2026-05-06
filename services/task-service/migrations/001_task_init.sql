BEGIN;

CREATE TABLE IF NOT EXISTS ai_models (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    token_cost NUMERIC(12,2) NOT NULL CHECK (token_cost >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS prompt_tasks (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    model_id TEXT NOT NULL REFERENCES ai_models(id) ON DELETE RESTRICT,
    payload TEXT NOT NULL,
    status TEXT NOT NULL,
    result TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_prompt_tasks_status
        CHECK (status IN ('Queued', 'Processing', 'Completed', 'Failed', 'Cancelled'))
);

CREATE TABLE IF NOT EXISTS worker_nodes (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_worker_nodes_status
        CHECK (status IN ('Idle', 'Busy'))
);

CREATE TABLE IF NOT EXISTS worker_supported_models (
    worker_id TEXT NOT NULL REFERENCES worker_nodes(id) ON DELETE CASCADE,
    model_id TEXT NOT NULL REFERENCES ai_models(id) ON DELETE CASCADE,
    PRIMARY KEY (worker_id, model_id)
);

CREATE INDEX IF NOT EXISTS idx_prompt_tasks_user_id ON prompt_tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_prompt_tasks_status ON prompt_tasks(status);
CREATE INDEX IF NOT EXISTS idx_prompt_tasks_created_at ON prompt_tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_worker_supported_models_model_id ON worker_supported_models(model_id);

INSERT INTO worker_nodes (id, status)
VALUES ('worker-1', 'Idle')
ON CONFLICT (id) DO NOTHING;

COMMIT;
