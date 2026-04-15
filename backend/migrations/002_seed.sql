-- Seed-дані для локальної розробки та демонстрацій.

BEGIN;

INSERT INTO users (id, username, token_balance) VALUES
    ('user-1', 'alice', 100.00),
    ('user-2', 'bob', 5.00),
    ('user-3', 'charlie', 200.00);

INSERT INTO ai_models (id, name, description, token_cost) VALUES
    ('model-1', 'Llama-3', 'Велика мовна модель для генерації тексту', 5.00),
    ('model-2', 'Stable-Diffusion', 'Модель для генерації зображень з тексту', 10.00),
    ('model-3', 'Whisper', 'Модель розпізнавання мовлення (Speech-to-text)', 3.00),
    ('model-4', 'GPT-4o', 'Просунута мультимодальна ШІ-модель', 15.00);

INSERT INTO worker_nodes (id, status) VALUES
    ('worker-1', 'Idle'),
    ('worker-2', 'Idle'),
    ('worker-3', 'Idle');

-- Кожен воркер підтримує всі доступні моделі.
INSERT INTO worker_supported_models (worker_id, model_id) VALUES
    ('worker-1', 'model-1'),
    ('worker-1', 'model-2'),
    ('worker-1', 'model-3'),
    ('worker-1', 'model-4'),
    ('worker-2', 'model-1'),
    ('worker-2', 'model-2'),
    ('worker-2', 'model-3'),
    ('worker-2', 'model-4'),
    ('worker-3', 'model-1'),
    ('worker-3', 'model-2'),
    ('worker-3', 'model-3'),
    ('worker-3', 'model-4');

COMMIT;
