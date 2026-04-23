CREATE TABLE sdk_keys
(
    id             UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    key_hash       VARCHAR(255) NOT NULL,
    project_id     UUID REFERENCES projects (id),
    environment_id UUID REFERENCES environments (id),
    name           VARCHAR(255) NOT NULL,
    created_by     UUID REFERENCES users (id),
    expires_at     TIMESTAMPTZ           DEFAULT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);