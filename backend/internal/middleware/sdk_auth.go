package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"net/http"
	"time"
)

type sdkContextKey string

const (
	SDKProjectIDKey     sdkContextKey = "sdk_project_id"
	SDKEnvironmentIDKey sdkContextKey = "sdk_environment_id"
)

type SDKKeyValidator interface {
	GetByKeyHash(ctx context.Context, hash string) (*domain.SDKKey, error)
}

type SDKAuthMiddleware struct {
	repo SDKKeyValidator
}

func NewSDKAuthMiddleware(repo SDKKeyValidator) *SDKAuthMiddleware {
	return &SDKAuthMiddleware{repo: repo}
}

func (m *SDKAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawKey := r.Header.Get("X-SDK-Key")
		if rawKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "missing X-SDK-Key header"}`))
			return
		}

		keyHash := hashSDKKey(rawKey)
		key, err := m.repo.GetByKeyHash(r.Context(), string(keyHash[:]))
		if err != nil || key == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid X-SDK-Key header"}`))
			return
		}
		if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "X-SDK-Key expired"}`))
			return
		}

		ctx := context.WithValue(r.Context(), SDKProjectIDKey, key.ProjectID)
		ctx = context.WithValue(ctx, SDKEnvironmentIDKey, key.EnvironmentID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hashSDKKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
