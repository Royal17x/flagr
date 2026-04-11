package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const flagCacheTTL = 5 * time.Minute

var ErrCacheMiss = errors.New("cache miss")

type FlagCache interface {
	GetEvaluation(ctx context.Context, flagKey string, projectID uuid.UUID, envID uuid.UUID) (bool, error)
	SetEvaluation(ctx context.Context, flagKey string, projectID uuid.UUID, envID uuid.UUID, enabled bool) error
	InvalidateFlag(ctx context.Context, flagKey string, projectID uuid.UUID) error
}

type redisCache struct {
	client *redis.Client
}

func NewFlagCache(client *redis.Client) FlagCache {
	return &redisCache{client: client}
}

func evalKey(flagKey string, projectID uuid.UUID, envID uuid.UUID) string {
	return fmt.Sprintf("flag:eval:%s:%s:%s", projectID, flagKey, envID)
}

func (r *redisCache) GetEvaluation(ctx context.Context, flagKey string, projectID uuid.UUID, envID uuid.UUID) (bool, error) {
	key := evalKey(flagKey, projectID, envID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, ErrCacheMiss
		}
		return false, fmt.Errorf("cache.GetEvaluation: %w", err)
	}
	res, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("failed to parse cache value: %w", err)
	}
	return res, nil
}

func (r *redisCache) SetEvaluation(ctx context.Context, flagKey string, projectID uuid.UUID, envID uuid.UUID, enabled bool) error {
	key := evalKey(flagKey, projectID, envID)
	err := r.client.Set(ctx, key, enabled, flagCacheTTL).Err()
	if err != nil {
		return fmt.Errorf("cache.SetEvaluation: %w", err)
	}
	return nil
}

func (r *redisCache) InvalidateFlag(ctx context.Context, flagKey string, projectID uuid.UUID) error {
	searchPattern := fmt.Sprintf("flag:eval:%s:%s:*", projectID, flagKey)
	// TODO: Change Keys with Scan for high load.
	keys, err := r.client.Keys(ctx, searchPattern).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("cache.GetEvaluation: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	err = r.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("cache.InvalidateFlag: %w", err)
	}
	return nil
}
