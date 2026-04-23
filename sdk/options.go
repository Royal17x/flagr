package sdk

import "time"

type Option func(*config)

func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

func WithGRPC() Option {
	return func(c *config) {
		c.useGRPC = true
	}
}

func WithCacheSize(size int) Option {
	return func(c *config) {
		c.cacheSize = size
	}
}

func WithCacheTTL(d time.Duration) Option {
	return func(c *config) {
		c.cacheTTL = d
	}
}

func WithDefaultValue(v bool) Option {
	return func(c *config) {
		c.defaultValue = v
	}
}

func WithSyncInterval(d time.Duration) Option {
	return func(c *config) {
		c.syncInterval = d
	}
}

func WithStreaming(enabled bool) Option {
	return func(c *config) {
		c.enableStreaming = enabled
	}
}

func WithTLS(enabled bool) Option {
	return func(c *config) {
		c.tlsEnabled = enabled
	}
}
