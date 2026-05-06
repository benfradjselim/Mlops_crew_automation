package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// operatorMetrics holds all instrumentation for the operator process.
// Plain int64 fields accessed via sync/atomic functions (compatible with Go 1.18).
type operatorMetrics struct {
	reconcileSuccess int64
	reconcileError   int64
	instancesCurrent int64
}

var metrics operatorMetrics

func recordReconcileSuccess() { atomic.AddInt64(&metrics.reconcileSuccess, 1) }
func recordReconcileError()   { atomic.AddInt64(&metrics.reconcileError, 1) }
func setInstanceCount(n int)  { atomic.StoreInt64(&metrics.instancesCurrent, int64(n)) }

// metricsHandler writes the Prometheus text exposition format.
func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	ok := atomic.LoadInt64(&metrics.reconcileSuccess)
	errCnt := atomic.LoadInt64(&metrics.reconcileError)
	instances := atomic.LoadInt64(&metrics.instancesCurrent)

	fmt.Fprintf(w, "# HELP ruptura_operator_reconcile_total Total reconcile calls by result.\n")
	fmt.Fprintf(w, "# TYPE ruptura_operator_reconcile_total counter\n")
	fmt.Fprintf(w, "ruptura_operator_reconcile_total{result=\"success\"} %d\n", ok)
	fmt.Fprintf(w, "ruptura_operator_reconcile_total{result=\"error\"} %d\n", errCnt)

	fmt.Fprintf(w, "# HELP ruptura_operator_instances_current Number of RupturaInstances currently managed.\n")
	fmt.Fprintf(w, "# TYPE ruptura_operator_instances_current gauge\n")
	fmt.Fprintf(w, "ruptura_operator_instances_current %d\n", instances)

	fmt.Fprintf(w, "# HELP ruptura_operator_info Operator version info.\n")
	fmt.Fprintf(w, "# TYPE ruptura_operator_info gauge\n")
	fmt.Fprintf(w, "ruptura_operator_info{version=%q} 1\n", operatorVersion)
}

// startMetricsServer starts a lightweight HTTP server on addr (e.g. ":9090").
// It blocks until the context is cancelled, then shuts down gracefully.
func startMetricsServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metricsHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logError("metrics server error", "err", err)
		}
	}()
	return srv
}
