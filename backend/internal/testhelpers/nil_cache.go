package testhelpers

import (
	"context"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/google/uuid"
)

// NilCache is a no-op cache implementation for use in tests.
type NilCache struct{}

func (n *NilCache) GetEvaluation(_ context.Context, _ string, _ uuid.UUID, _ uuid.UUID) (bool, error) {
	return false, cache.ErrCacheMiss
}
func (n *NilCache) SetEvaluation(_ context.Context, _ string, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return nil
}
func (n *NilCache) InvalidateFlag(_ context.Context, _ string, _ uuid.UUID) error {
	return nil
}
