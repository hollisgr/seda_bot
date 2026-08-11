-- +goose Up
-- +goose StatementBegin

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    tg_id INTEGER NOT NULL CHECK(tg_id > 0),
    chat_id INTEGER NOT NULL CHECK(chat_id > 0),
    name TEXT,
    first_name TEXT NOT NULL,
    last_name TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
