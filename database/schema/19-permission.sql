CREATE TABLE permission (
    id SERIAL PRIMARY KEY,
    principal_type INT NOT NULL,
    principal_value TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id INT NOT NULL,
    policy_id INT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL,
    update_time TIMESTAMPTZ NOT NULL,
    expiration_time TIMESTAMPTZ NOT NULL,
    removed BOOLEAN NOT NULL,
    UNIQUE (
        principal_type,
        principal_value,
        resource_type,
        resource_id
    )
);
CREATE INDEX idx_policy_id ON permission (policy_id);