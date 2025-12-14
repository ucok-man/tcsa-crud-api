-- +goose Up
-- +goose StatementBegin
INSERT INTO "transactions" (user_id, amount, status, created_at, updated_at)
SELECT 
    -- Random user_id between 1 and 10
    (random() * 9 + 1)::BIGINT,
    -- Random amount between 10000 and 5000000 (100 to 50,000 if in cents)
    (random() * 4990000 + 10000)::BIGINT,
    -- Random status
    CASE 
        WHEN random() < 0.7 THEN 'success'::transaction_status
        WHEN random() < 0.9 THEN 'pending'::transaction_status
        ELSE 'failed'::transaction_status
    END,
    -- Random timestamp between 1 month ago and now
    CURRENT_TIMESTAMP - (random() * INTERVAL '30 days'),
    -- updated_at is current time
    CURRENT_TIMESTAMP
FROM generate_series(1, 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "transactions"
WHERE user_id BETWEEN 1 AND 10;
-- +goose StatementEnd