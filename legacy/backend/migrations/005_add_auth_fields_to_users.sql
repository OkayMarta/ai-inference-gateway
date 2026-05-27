BEGIN;

ALTER TABLE users
    ADD COLUMN email TEXT,
    ADD COLUMN password_hash TEXT,
    ADD COLUMN role TEXT NOT NULL DEFAULT 'user',
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE users
    ADD CONSTRAINT users_email_unique UNIQUE (email);

DELETE FROM transactions;
DELETE FROM prompt_tasks;
DELETE FROM users;

ALTER TABLE users
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password_hash SET NOT NULL;

COMMIT;
