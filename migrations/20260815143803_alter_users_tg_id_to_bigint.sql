-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN tg_id TYPE BIGINT;
ALTER TABLE users ALTER COLUMN chat_id TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN tg_id TYPE INT;
ALTER TABLE users ALTER COLUMN chat_id TYPE INT;
-- +goose StatementEnd
