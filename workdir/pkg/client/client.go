package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/benfradjselim/ruptura/pkg/models"
)

// Client is a typed HTTP client for the Ruptura v2 REST API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config holds client configuration.
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

func New(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	return &Client{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(raw))
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doRaw(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// --- Response types ---

type HealthResponse struct {
	Status           string `json:"status"`
	Version          string `json:"version"`
	RuptureDetection string `json:"rupture_detection"`
	Message          string `json:"message"`
	UptimeSeconds    int64  `json:"uptime_seconds"`
	Edition          string `json:"edition"`
}

type ActionItem struct {
	ID         string    `json:"id"`
	WorkloadID string    `json:"workload_id"`
	Type       string    `json:"type"`
	Tier       int       `json:"tier"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	FusedR     float64   `json:"fused_rupture_index"`
}

type Suppression struct {
	ID       string    `json:"id"`
	Workload string    `json:"workload"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Reason   string    `json:"reason"`
}

type CreateSuppressionReq struct {
	Workload string    `json:"workload"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Reason   string    `json:"reason"`
}

type SimInjectReq struct {
	Pattern         string `json:"pattern"`
	Workload        string `json:"workload,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
}

type SimInjectResp struct {
	Pattern  string `json:"pattern"`
	Workload string `json:"workload"`
	Message  string `json:"message"`
}

type ExplainNarrative struct {
	RuptureID string `json:"rupture_id"`
	Narrative string `json:"narrative"`
	Summary   string `json:"summary"`
}

type AnomalyEvent struct {
	Host     string `json:"host"`
	Severity string `json:"severity"`
	Total    int64  `json:"total"`
}

type ContextEntry struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Service   string    `json:"service"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// --- API methods ---

func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	var r HealthResponse
	return r, c.do(ctx, http.MethodGet, "/api/v2/health", nil, &r)
}

func (c *Client) Metrics(ctx context.Context) (string, error) {
	b, err := c.doRaw(ctx, "/api/v2/metrics")
	return string(b), err
}

func (c *Client) Snapshots(ctx context.Context) ([]models.KPISnapshot, error) {
	var r []models.KPISnapshot
	return r, c.do(ctx, http.MethodGet, "/api/v2/ruptures", nil, &r)
}

func (c *Client) Snapshot(ctx context.Context, ref string) (models.KPISnapshot, error) {
	var r models.KPISnapshot
	return r, c.do(ctx, http.MethodGet, "/api/v2/rupture/"+ref, nil, &r)
}

func (c *Client) Actions(ctx context.Context) ([]ActionItem, error) {
	var r []ActionItem
	err := c.do(ctx, http.MethodGet, "/api/v2/actions", nil, &r)
	if r == nil {
		r = []ActionItem{}
	}
	return r, err
}

func (c *Client) ApproveAction(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/api/v2/actions/"+id+"/approve", nil, nil)
}

func (c *Client) RejectAction(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/api/v2/actions/"+id+"/reject", nil, nil)
}

func (c *Client) EmergencyStop(ctx context.Context) error {
	return c.do(ctx, http.MethodPost, "/api/v2/actions/emergency-stop", nil, nil)
}

func (c *Client) Suppressions(ctx context.Context) ([]Suppression, error) {
	var r []Suppression
	err := c.do(ctx, http.MethodGet, "/api/v2/suppressions", nil, &r)
	if r == nil {
		r = []Suppression{}
	}
	return r, err
}

func (c *Client) CreateSuppression(ctx context.Context, req CreateSuppressionReq) (Suppression, error) {
	var r Suppression
	return r, c.do(ctx, http.MethodPost, "/api/v2/suppressions", req, &r)
}

func (c *Client) DeleteSuppression(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/api/v2/suppressions/"+id, nil, nil)
}

func (c *Client) Weights(ctx context.Context) ([]models.SignalWeights, error) {
	var r []models.SignalWeights
	err := c.do(ctx, http.MethodGet, "/api/v2/config/weights", nil, &r)
	if r == nil {
		r = []models.SignalWeights{}
	}
	return r, err
}

func (c *Client) SetWeights(ctx context.Context, cfgs []models.SignalWeights) error {
	return c.do(ctx, http.MethodPost, "/api/v2/config/weights", cfgs, nil)
}

func (c *Client) SimInject(ctx context.Context, req SimInjectReq) (SimInjectResp, error) {
	var r SimInjectResp
	return r, c.do(ctx, http.MethodPost, "/api/v2/sim/inject", req, &r)
}

func (c *Client) Explain(ctx context.Context, ruptureID string) (ExplainNarrative, error) {
	var r ExplainNarrative
	return r, c.do(ctx, http.MethodGet, "/api/v2/explain/"+ruptureID+"/narrative", nil, &r)
}

func (c *Client) Anomalies(ctx context.Context) ([]AnomalyEvent, error) {
	var r []AnomalyEvent
	err := c.do(ctx, http.MethodGet, "/api/v2/anomalies", nil, &r)
	if r == nil {
		r = []AnomalyEvent{}
	}
	return r, err
}

func (c *Client) AddContext(ctx context.Context, entry ContextEntry) (ContextEntry, error) {
	var r ContextEntry
	return r, c.do(ctx, http.MethodPost, "/api/v2/context", entry, &r)
}

func (c *Client) ListContexts(ctx context.Context) ([]ContextEntry, error) {
	var r []ContextEntry
	return r, c.do(ctx, http.MethodGet, "/api/v2/context", nil, &r)
}

func (c *Client) DeleteContext(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/api/v2/context/"+id, nil, nil)
}
