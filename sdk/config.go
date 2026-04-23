package sdk

import (
	"time"
)

type config struct {
	// flagr server address
	serverURL string
	// sdk key for auth
	sdkKey string
	// use gRPC instead of HTTP (faster)
	useGRPC bool
	// timeout for request
	timeout time.Duration
	// size of local cache (amount of flags)
	cacheSize int
	// ttl of local cache
	cacheTTL time.Duration
	// if flagr is unavailable - def value
	defaultValue bool
	// background flag synchronization interval
	syncInterval time.Duration
	// enable WatchFlags streaming (auto-update)
	enableStreaming bool
	// TLS for gRPC connection
	tlsEnabled bool
}

func defaultConfig() *config {
	return &config{
		timeout:         time.Second * 5,
		cacheSize:       10000,
		cacheTTL:        time.Minute * 5,
		defaultValue:    false,
		syncInterval:    time.Second * 30,
		enableStreaming: true,
		tlsEnabled:      true,
	}
}
