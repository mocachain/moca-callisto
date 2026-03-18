-- buckets tables
CREATE TABLE buckets (
    id SERIAL PRIMARY KEY,
    bucket_id TEXT NOT NULL UNIQUE,
    bucket_name TEXT NOT NULL CHECK (
        LENGTH(bucket_name) BETWEEN 3 AND 63
    ),
    owner_address TEXT NOT NULL,
    payment_address TEXT,
    global_virtual_group_family_id INT NOT NULL,
    operator_address TEXT,
    source_type TEXT,
    charged_read_quota BIGINT,
    visibility TEXT,
    status TEXT,
    delete_at BIGINT,
    delete_reason TEXT,
    storage_size DECIMAL(65, 0) NOT NULL,
    charge_size DECIMAL(65, 0) NOT NULL,
    create_at BIGINT,
    create_tx_hash TEXT NOT NULL,
    create_evm_tx_hash TEXT NOT NULL,
    create_time TIMESTAMPTZ,
    update_at BIGINT,
    update_tx_hash TEXT NOT NULL,
    update_evm_tx_hash TEXT NOT NULL,
    update_time TIMESTAMPTZ,
    removed BOOLEAN DEFAULT FALSE,
    off_chain_status INT NOT NULL DEFAULT 0,
    sp_as_delegated_agent_disabled BOOLEAN DEFAULT FALSE,
    tags JSONB,
    UNIQUE (bucket_id),
    CONSTRAINT fk_family_id FOREIGN KEY (global_virtual_group_family_id) REFERENCES global_virtual_group_families(global_virtual_group_family_id)
);
CREATE INDEX idx_bucket_owner ON buckets(owner_address);
CREATE INDEX idx_bucket_gvgf_id ON buckets(global_virtual_group_family_id);