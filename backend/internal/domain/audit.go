package domain

import (
	"github.com/google/uuid"
	"time"
)

type AuditEntry struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Action     string    `json:"action" db:"action"`
	ActorID    uuid.UUID `json:"actor_id" db:"actor_id"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
	Payload    string    `json:"payload" db:"payload"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
