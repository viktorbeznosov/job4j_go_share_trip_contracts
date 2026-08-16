-- +goose Up
-- +goose StatementBegin
CREATE TABLE contracts_services (
    contract_id uuid,
    service_id uuid
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contracts_services;
-- +goose StatementEnd
