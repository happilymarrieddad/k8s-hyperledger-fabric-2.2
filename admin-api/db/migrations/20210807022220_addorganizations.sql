-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizations(
    id bigserial NOT NULL PRIMARY KEY,
    name citext NOT NULL,
    network_name citext NOT NULL,
    namespace_created bool DEFAULT FALSE,
    active bool DEFAULT TRUE,
    number_of_cas bigint NOT NULL,
    number_of_peers bigint NOT NULL,
    network_active bool DEFAULT FALSE,
    UNIQUE(name),
    UNIQUE(network_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd
