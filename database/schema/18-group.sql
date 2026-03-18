CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    owner_address TEXT NOT NULL,
    group_id TEXT NOT NULL,
    group_name TEXT,
    source_type TEXT,
    extra TEXT,
    account_address TEXT NOT NULL,
    operator_address TEXT,
    expiration_time TIMESTAMPTZ,
    create_at BIGINT,
    create_time TIMESTAMPTZ,
    create_tx_hash TEXT NOT NULL,
    create_evm_tx_hash TEXT NOT NULL,
    update_at BIGINT,
    update_time TIMESTAMPTZ,
    update_tx_hash TEXT NOT NULL,
    update_evm_tx_hash TEXT NOT NULL,
    removed BOOLEAN DEFAULT FALSE,
    tags JSONB,
    CONSTRAINT idx_account_group UNIQUE (account_address, group_id)
);
CREATE INDEX idx_owner ON groups(owner_address);
CREATE INDEX idx_group_id ON groups(group_id);
CREATE INDEX idx_group_name ON groups(group_name);