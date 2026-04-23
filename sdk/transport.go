package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Evaluator interface {
	Evaluate(ctx context.Context, req EvaluateRequest) (EvaluateResponse, error)
	Close() error
}

type EvaluateRequest struct {
	FlagKey       string
	ProjectID     string
	EnvironmentID string
	Context       map[string]string
}
type EvaluateResponse struct {
	Enabled          bool
	Reason           string
	EvaluationTimeMs int64
}

type httpEvaluator struct {
	baseURL string
	sdkKey  string
	client  *http.Client
}

func newHTTPEvaluator(baseURL, sdkKey string, timeout time.Duration) *httpEvaluator {
	return &httpEvaluator{
		baseURL: baseURL,
		sdkKey:  sdkKey,
		client:  &http.Client{Timeout: timeout},
	}
}

func (e *httpEvaluator) Evaluate(ctx context.Context, req EvaluateRequest) (EvaluateResponse, error) {
	url := fmt.Sprintf("%s/api/v1/flags/evaluate?key=%s&project_id=%s&environment_id=%s",
		e.baseURL, req.FlagKey, req.ProjectID, req.EnvironmentID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return EvaluateResponse{}, fmt.Errorf("evaluate: create request: %w", err)
	}

	httpReq.Header.Set("X-SDK-Key", e.sdkKey)

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return EvaluateResponse{}, fmt.Errorf("evaluate: do request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return EvaluateResponse{}, fmt.Errorf("evaluate: decode: %w", err)
	}

	return EvaluateResponse{
		Enabled: result.Enabled,
		Reason:  "HTTP",
	}, nil
}

func (e *httpEvaluator) Close() error { return nil }
