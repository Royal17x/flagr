package domain

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Project struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Environment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type FlagType string

const (
	FlagTypeBoolean      FlagType = "boolean"
	FlagTypeMultivariate FlagType = "multivariate"
)

type Flag struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Key         string    `json:"key" db:"key"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        FlagType  `json:"type" db:"type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type FlagEnvironment struct {
	ID                uuid.UUID `json:"id" db:"id"`
	FlagID            uuid.UUID `json:"flag_id" db:"flag_id"`
	EnvironmentID     uuid.UUID `json:"environment_id" db:"environment_id"`
	Enabled           bool      `json:"enabled" db:"enabled"`
	RolloutPercentage int       `json:"rollout_percentage" db:"rollout_percentage"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy         uuid.UUID `json:"updated_by" db:"updated_by"`
}

type Rule struct {
	ID                uuid.UUID `json:"id" db:"id"`
	FlagEnvironmentID uuid.UUID `json:"flag_environment_id" db:"flag_environment_id"`
	Attribute         string    `json:"attribute" db:"attribute"`
	Operator          string    `json:"operator" db:"operator"`
	Value             string    `json:"value" db:"value"`
	Priority          int       `json:"priority" db:"priority"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type AuditEntry struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Action     string    `json:"action" db:"action"`
	ActorID    uuid.UUID `json:"actor_id" db:"actor_id"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
	Payload    string    `json:"payload" db:"payload"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
