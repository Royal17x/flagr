package sdk

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Client struct {
	cfg       *config
	evaluator Evaluator
	cache     *localCache

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewClient(serverURL, sdkKey string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	cfg.serverURL = serverURL
	cfg.sdkKey = sdkKey

	for _, opt := range opts {
		opt(cfg)
	}
	var evaluator Evaluator
	var err error
	if cfg.useGRPC {
		evaluator, err = newGRPCEvaluator(serverURL, sdkKey, cfg.timeout, cfg.tlsEnabled)
		if err != nil {
			return nil, fmt.Errorf("flagr: create grpc evaluator: %w", err)
		}
	} else {
		evaluator = newHTTPEvaluator(serverURL, sdkKey, cfg.timeout)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		cfg:       cfg,
		evaluator: evaluator,
		cache:     newLocalCache(cfg.cacheSize, cfg.cacheTTL),
		ctx:       ctx,
		cancel:    cancel,
	}

	return c, nil
}

func (c *Client) IsEnabled(ctx context.Context, flagKey, projectID, environmentID string) bool {
	if enabled, ok := c.cache.get(flagKey, projectID, environmentID); ok {
		return enabled
	}

	resp, err := c.evaluator.Evaluate(ctx, EvaluateRequest{
		FlagKey:       flagKey,
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	})
	if err != nil {
		slog.Warn("flagr: evaluate failed, using default",
			"flag", flagKey,
			"default", c.cfg.defaultValue,
			"error", err,
		)
		return c.cfg.defaultValue
	}

	c.cache.set(flagKey, projectID, environmentID, resp.Enabled)
	return resp.Enabled
}

func (c *Client) IsEnabledWithContext(
	ctx context.Context,
	flagKey, projectID, environmentID string,
	userContext map[string]string,
) bool {
	// TODO: Targetting rules
	return c.IsEnabled(ctx, flagKey, projectID, environmentID)
}

func (c *Client) startBackgroundSync(projectID, environmentID string) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.cfg.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				c.cache.invalidateAll()
				slog.Debug("flagr: cache invalidated by background sync")
			}
		}
	}()
}

func (c *Client) Close() error {
	c.cancel()
	c.wg.Wait()
	return c.evaluator.Close()
}
