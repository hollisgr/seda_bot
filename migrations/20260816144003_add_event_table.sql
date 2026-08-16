-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    id SERIAL PRIMARY KEY,
    type TEXT,
    name TEXT,
    description TEXT,
    date TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
