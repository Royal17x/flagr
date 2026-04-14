package cache

import (
	"context"

	"github.com/google/uuid"
)

type NilCache struct{}

func (n *NilCache) GetEvaluation(_ context.Context, _ string, _ uuid.UUID, _ uuid.UUID) (bool, error) {
	return false, ErrCacheMiss
}
func (n *NilCache) SetEvaluation(_ context.Context, _ string, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return nil
}
func (n *NilCache) InvalidateFlag(_ context.Context, _ string, _ uuid.UUID) error {
	return nil
}
