package testhelpers

import (
	"context"
	"github.com/Royal17x/flagr/backend/pkg/kafka"
)

type NilAuditPublisher struct{}

func (n *NilAuditPublisher) PublishAuditEvent(_ context.Context, _ kafka.AuditMessage) error {
	return nil
}
