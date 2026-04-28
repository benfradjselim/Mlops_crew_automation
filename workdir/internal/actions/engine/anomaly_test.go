package engine

import (
	"testing"
	"time"

	"github.com/benfradjselim/ruptura/pkg/models"
)

func TestRecommendFromAnomaly_criticalFiresRecommendations(t *testing.T) {
	e, _ := New(nil, nil)
	ev := models.AnomalyEvent{
		Host:      "db-primary",
		Metric:    "cpu_percent",
		Score:     7.5,
		Severity:  models.SeverityCritical,
		Timestamp: time.Now(),
	}

	recs, err := e.RecommendFromAnomaly(ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Error("expected non-empty recommendations for SeverityCritical anomaly")
	}
	for _, r := range recs {
		if r.Host != "db-primary" {
			t.Errorf("expected host db-primary, got %q", r.Host)
		}
	}
}

func TestRecommendFromAnomaly_nonCriticalReturnsNil(t *testing.T) {
	e, _ := New(nil, nil)
	ev := models.AnomalyEvent{
		Host:      "web-1",
		Metric:    "latency",
		Score:     2.0,
		Severity:  models.SeverityWarning,
		Timestamp: time.Now(),
	}

	recs, err := e.RecommendFromAnomaly(ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected no recommendations for non-critical anomaly, got %d", len(recs))
	}
}

func TestRecommendFromAnomaly_profileAssignment(t *testing.T) {
	e, _ := New(nil, nil)

	// Score > 5.0 → spike profile, matching default-spike (MinR=3.0, Profile=spike)
	// and default-any (MinR=5.0, Profile="")
	ev := models.AnomalyEvent{
		Host:      "host-x",
		Metric:    "cpu",
		Score:     6.0,
		Severity:  models.SeverityCritical,
		Timestamp: time.Now(),
	}
	recs, err := e.RecommendFromAnomaly(ev)
	if err != nil {
		t.Fatal(err)
	}
	// Should have spike + any match at minimum
	if len(recs) < 2 {
		t.Errorf("expected at least 2 recommendations for score>5 critical, got %d", len(recs))
	}
}
