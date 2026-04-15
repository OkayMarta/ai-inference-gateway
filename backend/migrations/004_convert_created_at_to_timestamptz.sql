BEGIN;

ALTER TABLE prompt_tasks
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING created_at AT TIME ZONE current_setting('TIMEZONE');

ALTER TABLE transactions
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING created_at AT TIME ZONE current_setting('TIMEZONE');

ALTER TABLE worker_nodes
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING created_at AT TIME ZONE current_setting('TIMEZONE');

COMMIT;
