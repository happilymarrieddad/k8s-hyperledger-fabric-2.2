-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizations(
    id bigserial NOT NULL PRIMARY KEY,
    name citext NOT NULL,
    network_name citext NOT NULL,
    namespace_created bool DEFAULT FALSE,
    active bool DEFAULT TRUE,
    UNIQUE(name),
    UNIQUE(network_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd
