package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew_DefaultTimeout(t *testing.T) {
	c := New(Config{BaseURL: "http://localhost"})
	if c.httpClient.Timeout != 15*time.Second {
		t.Errorf("want 15s, got %v", c.httpClient.Timeout)
	}
}

func TestHealth_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ready","version":"6.6.0"}`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Health(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "ready" {
		t.Errorf("want ready, got %s", res.Status)
	}
	if res.Version != "6.6.0" {
		t.Errorf("want 6.6.0, got %s", res.Version)
	}
}

func TestHealth_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	_, err := c.Health(context.Background())
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Errorf("expected 500 error, got %v", err)
	}
}

func TestSnapshots_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Snapshots(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("want empty slice, got %v", res)
	}
}

func TestSnapshot_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"host":"payment-api","fused_rupture_index":2.1}`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Snapshot(context.Background(), "payment-api")
	if err != nil {
		t.Fatal(err)
	}
	if res.Host != "payment-api" {
		t.Errorf("want payment-api, got %s", res.Host)
	}
	if res.FusedRuptureIndex != 2.1 {
		t.Errorf("want 2.1, got %f", res.FusedRuptureIndex)
	}
}

func TestActions_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"act_1","type":"scale","tier":2}]`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Actions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].ID != "act_1" {
		t.Errorf("unexpected response: %+v", res)
	}
}

func TestApproveAction_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	if err := c.ApproveAction(context.Background(), "act_1"); err != nil {
		t.Fatal(err)
	}
}

func TestEmergencyStop_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	if err := c.EmergencyStop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWeights_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"selector":"production/*","stress":0.4}]`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Weights(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].Selector != "production/*" {
		t.Errorf("unexpected: %+v", res)
	}
}

func TestAddContext_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"ctx_1"}`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.AddContext(context.Background(), ContextEntry{ID: "ctx_1"})
	if err != nil {
		t.Fatal(err)
	}
	if res.ID != "ctx_1" {
		t.Errorf("want ctx_1, got %s", res.ID)
	}
}

func TestMetrics_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "rpt_uptime_seconds 120")
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL})
	res, err := c.Metrics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res, "rpt_uptime_seconds") {
		t.Errorf("unexpected metrics: %s", res)
	}
}

func TestAuth_HeaderSent(t *testing.T) {
	token := "secret-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+token {
			t.Errorf("want Bearer %s, got %s", token, auth)
		}
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer ts.Close()

	c := New(Config{BaseURL: ts.URL, APIKey: token})
	c.Health(context.Background()) //nolint:errcheck
}

func TestClient_NetworkError(t *testing.T) {
	c := New(Config{BaseURL: "http://127.0.0.1:19999"})
	_, err := c.Health(context.Background())
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}
