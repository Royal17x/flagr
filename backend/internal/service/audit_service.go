package service

import (
	"context"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/pkg/kafka"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type AuditPublisher interface {
	PublishAuditEvent(ctx context.Context, event kafka.AuditMessage) error
}

type AuditService struct {
	publisher AuditPublisher
}

func NewAuditService(publisher AuditPublisher) *AuditService {
	return &AuditService{publisher: publisher}
}

func (a *AuditService) publish(ctx context.Context, msg kafka.AuditMessage) {
	if err := a.publisher.PublishAuditEvent(ctx, msg); err != nil {
		slog.Error("audit: failed to publish event",
			"action", msg.Action,
			"resource_id", msg.ResourceID,
			"actor_id", msg.ActorID,
			"error", err,
		)
	}
}

func (a *AuditService) RecordFlagCreated(ctx context.Context, actorID, flagID, projectID uuid.UUID) {
	a.publish(ctx, kafka.AuditMessage{
		ID:         uuid.New(),
		Version:    1,
		Action:     string(domain.AuditActionFlagCreated),
		ActorID:    actorID,
		ResourceID: flagID,
		ProjectID:  projectID,
		OccurredAt: time.Now().UTC(),
	})
}

func (a *AuditService) RecordFlagUpdated(ctx context.Context, actorID, flagID, projectID uuid.UUID, changes any) {
	a.publish(ctx, kafka.AuditMessage{
		ID:         uuid.New(),
		Version:    1,
		Action:     string(domain.AuditActionFlagUpdated),
		ActorID:    actorID,
		ResourceID: flagID,
		ProjectID:  projectID,
		Payload:    changes,
		OccurredAt: time.Now().UTC(),
	})
}

func (a *AuditService) RecordFlagDeleted(ctx context.Context, actorID, flagID, projectID uuid.UUID) {
	a.publish(ctx, kafka.AuditMessage{
		ID:         uuid.New(),
		Version:    1,
		Action:     string(domain.AuditActionFlagDeleted),
		ActorID:    actorID,
		ResourceID: flagID,
		ProjectID:  projectID,
		OccurredAt: time.Now().UTC(),
	})
}
