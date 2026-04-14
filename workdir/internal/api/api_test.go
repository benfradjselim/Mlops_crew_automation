package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benfradjselim/ohe/internal/alerter"
	"github.com/benfradjselim/ohe/internal/analyzer"
	"github.com/benfradjselim/ohe/internal/api"
	"github.com/benfradjselim/ohe/internal/predictor"
	"github.com/benfradjselim/ohe/internal/processor"
	"github.com/benfradjselim/ohe/internal/storage"
	"github.com/benfradjselim/ohe/pkg/models"
)

func setupServer(t *testing.T) *httptest.Server {
	t.Helper()
	dir, err := os.MkdirTemp("", "ohe-api-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	store, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("Open storage: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	proc := processor.NewProcessor(1000)
	ana := analyzer.NewAnalyzer()
	pred := predictor.NewPredictor()
	alrt := alerter.NewAlerter(100)

	handlers := api.NewHandlers(store, proc, ana, pred, alrt, "test-secret", false)
	router := api.NewRouter(handlers, "test-secret", false)
	return httptest.NewServer(router)
}

func TestHealthEndpoint(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["success"] != true {
		t.Error("success should be true")
	}
}

func TestIngestAndKPIs(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	batch := models.MetricBatch{
		AgentID:   "test-agent",
		Host:      "testhost",
		Timestamp: time.Now(),
		Metrics: []models.Metric{
			{Name: "cpu_percent", Value: 60, Host: "testhost", Timestamp: time.Now()},
			{Name: "memory_percent", Value: 70, Host: "testhost", Timestamp: time.Now()},
			{Name: "load_avg_1", Value: 1.5, Host: "testhost", Timestamp: time.Now()},
		},
	}
	body, _ := json.Marshal(batch)
	resp, err := http.Post(srv.URL+"/api/v1/ingest", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /ingest: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("ingest status = %d; want 200", resp.StatusCode)
	}

	// Check KPIs are computed and accessible
	resp2, err := http.Get(srv.URL + "/api/v1/kpis?host=testhost")
	if err != nil {
		t.Fatalf("GET /kpis: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("kpis status = %d; want 200", resp2.StatusCode)
	}
}

func TestMetricsListEndpoint(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/metrics?host=localhost")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestDashboardCRUD(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	// Create
	d := models.Dashboard{Name: "Test Dashboard", Refresh: 30}
	body, _ := json.Marshal(d)
	resp, err := http.Post(srv.URL+"/api/v1/dashboards", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /dashboards: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("create status = %d; want 201", resp.StatusCode)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	data := created["data"].(map[string]interface{})
	id := data["id"].(string)

	// Get
	resp2, err := http.Get(srv.URL + "/api/v1/dashboards/" + id)
	if err != nil {
		t.Fatalf("GET /dashboards/%s: %v", id, err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("get status = %d; want 200", resp2.StatusCode)
	}

	// List
	resp3, err := http.Get(srv.URL + "/api/v1/dashboards")
	if err != nil {
		t.Fatalf("GET /dashboards: %v", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("list status = %d; want 200", resp3.StatusCode)
	}

	// Delete
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/api/v1/dashboards/"+id, nil)
	resp4, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /dashboards/%s: %v", id, err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d; want 204", resp4.StatusCode)
	}
}

func TestAlertsEndpoint(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/alerts")
	if err != nil {
		t.Fatalf("GET /alerts: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestQueryEndpoint(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	qr := models.QueryRequest{
		Query: "cpu_percent",
		From:  time.Now().Add(-time.Hour),
		To:    time.Now(),
	}
	body, _ := json.Marshal(qr)
	resp, err := http.Post(srv.URL+"/api/v1/query?host=localhost", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /query: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("query status = %d; want 200", resp.StatusCode)
	}
}

func TestPredictEndpoint(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/predict?host=localhost&horizon=60")
	if err != nil {
		t.Fatalf("GET /predict: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestCORSHeaders(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodOptions, srv.URL+"/api/v1/health", nil)
	req.Header.Set("Origin", "https://example.com")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS: %v", err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Access-Control-Allow-Origin") == "" {
		t.Error("missing CORS header")
	}
}
