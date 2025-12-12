-- +goose Up
-- +goose StatementBegin
CREATE DOMAIN "transaction_status" AS TEXT
    CONSTRAINT "valid_transaction_status" CHECK (VALUE IN ('success', 'pending', 'failed'));

CREATE TABLE IF NOT EXISTS "transactions" (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "transactions";
DROP DOMAIN IF EXISTS "transaction_status";
-- +goose StatementEnd
