-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN state TEXT NOT NULL DEFAULT 'main_menu';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN state;
-- +goose StatementEnd