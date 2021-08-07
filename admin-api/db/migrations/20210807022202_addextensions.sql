-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION citext;
-- +goose StatementEnd
