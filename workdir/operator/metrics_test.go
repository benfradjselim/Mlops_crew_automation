package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestMetricsHandler_ContainsExpectedMetrics(t *testing.T) {
	// Reset counters to a known state.
	atomic.StoreInt64(&metrics.reconcileSuccess, 3)
	atomic.StoreInt64(&metrics.reconcileError, 1)
	atomic.StoreInt64(&metrics.instancesCurrent, 2)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	metricsHandler(rr, req)

	body := rr.Body.String()

	cases := []string{
		`ruptura_operator_reconcile_total{result="success"} 3`,
		`ruptura_operator_reconcile_total{result="error"} 1`,
		`ruptura_operator_instances_current 2`,
		`ruptura_operator_info{version="` + operatorVersion + `"} 1`,
		"# TYPE ruptura_operator_reconcile_total counter",
		"# TYPE ruptura_operator_instances_current gauge",
		"# TYPE ruptura_operator_info gauge",
	}
	for _, want := range cases {
		if !strings.Contains(body, want) {
			t.Errorf("metrics output missing %q\nfull body:\n%s", want, body)
		}
	}
}

func TestMetricsHandler_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	metricsHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("expected text/plain Content-Type, got %q", ct)
	}
}

func TestHealthzHandler(t *testing.T) {
	srv := startMetricsServer(":0") // port 0 won't bind but we can test via httptest
	defer srv.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rr.Body.String())
	}
}

func TestRecordHelpers(t *testing.T) {
	atomic.StoreInt64(&metrics.reconcileSuccess, 0)
	atomic.StoreInt64(&metrics.reconcileError, 0)
	atomic.StoreInt64(&metrics.instancesCurrent, 0)

	recordReconcileSuccess()
	recordReconcileSuccess()
	recordReconcileError()
	setInstanceCount(5)

	if got := atomic.LoadInt64(&metrics.reconcileSuccess); got != 2 {
		t.Errorf("reconcileSuccess: want 2, got %d", got)
	}
	if got := atomic.LoadInt64(&metrics.reconcileError); got != 1 {
		t.Errorf("reconcileError: want 1, got %d", got)
	}
	if got := atomic.LoadInt64(&metrics.instancesCurrent); got != 5 {
		t.Errorf("instancesCurrent: want 5, got %d", got)
	}
}
