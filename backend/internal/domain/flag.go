package domain

import (
	"github.com/google/uuid"
	"time"
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

type FlagType string

const (
	FlagTypeBoolean      FlagType = "boolean"
	FlagTypeMultivariate FlagType = "multivariate"
)

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
