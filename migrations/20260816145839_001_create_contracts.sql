-- +goose Up
-- +goose StatementBegin

CREATE TYPE contract_status AS ENUM (
    'draft',
    'active',
    'suspended',
    'terminated'
);

CREATE TABLE contracts (
    id UUID PRIMARY KEY,
    company_id UUID,
    user_id UUID,
    status contract_status NOT NULL DEFAULT 'draft',
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contracts;
-- +goose StatementEnd
