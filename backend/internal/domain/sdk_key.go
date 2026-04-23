package domain

import (
	"github.com/google/uuid"
	"time"
)

type SDKKey struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	KeyHash       string     `json:"-" db:"key_hash"`
	ProjectID     uuid.UUID  `json:"project_id" db:"project_id"`
	EnvironmentID uuid.UUID  `json:"environment_id" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	CreatedBy     uuid.UUID  `json:"created_by" db:"created_by"`
	ExpiresAt     *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}
