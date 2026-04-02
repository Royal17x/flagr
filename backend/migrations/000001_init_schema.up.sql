CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE organizations
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE projects
(
    id              UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations (id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL UNIQUE,
    description     VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_org_id ON projects (organization_id);

CREATE TABLE environments
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects (id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE flags
(
    id          UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    project_id  UUID REFERENCES projects (id) ON DELETE CASCADE,
    key         VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    type        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_flags_project_id ON flags (project_id);

CREATE TABLE flag_environments
(
    id                 UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    flag_id            UUID        NOT NULL REFERENCES flags (id) ON DELETE CASCADE,
    environment_id     UUID        NOT NULL REFERENCES environments (id) ON DELETE CASCADE,
    enabled            BOOLEAN     NOT NULL DEFAULT false,
    rollout_percentage SMALLINT             DEFAULT 0 CHECK (rollout_percentage BETWEEN 0 AND 100),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by         UUID,
    UNIQUE (flag_id, environment_id)
);

CREATE INDEX idx_flag_env_flag_id ON flag_environments (flag_id);
CREATE INDEX idx_flag_env_environment_id ON flag_environments (environment_id);

CREATE TABLE rules
(
    id                  UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    flag_environment_id UUID         NOT NULL REFERENCES flag_environments (id) ON DELETE CASCADE,
    attribute           VARCHAR(100) NOT NULL,
    operator            VARCHAR(20)  NOT NULL,
    value               VARCHAR(100) NOT NULL,
    priority            INTEGER      NOT NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rules_flag_environment_id ON rules (flag_environment_id);
CREATE INDEX idx_rules_priority ON rules (priority);

CREATE TABLE audit_entries
(
    id          UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    action      VARCHAR(255) NOT NULL,
    actor_id    UUID         NOT NULL,
    resource_id UUID         NOT NULL,
    payload     TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);