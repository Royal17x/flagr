CREATE TABLE users
(
    id            UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    org_id        UUID REFERENCES organizations (id) ON DELETE CASCADE,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);