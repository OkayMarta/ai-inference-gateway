-- Seed-дані для локальної розробки та демонстрацій.

BEGIN;

INSERT INTO users (id, username, token_balance) VALUES
    ('user-1', 'alice', 100.00),
    ('user-2', 'bob', 5.00),
    ('user-3', 'charlie', 200.00);

INSERT INTO worker_nodes (id, status) VALUES
    ('worker-1', 'Idle'),
    ('worker-2', 'Idle'),
    ('worker-3', 'Idle');

COMMIT;
