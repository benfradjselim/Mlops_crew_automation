package storage

import (
	"os"
	"testing"
	"time"
)

func TestStorageOpenClose(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !s.Healthy() {
		t.Error("store should be healthy")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestSaveAndGetMetric(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	now := time.Now()
	if err := s.SaveMetric("host1", "cpu_percent", 0.75, now); err != nil {
		t.Fatalf("SaveMetric: %v", err)
	}

	from := now.Add(-time.Second)
	to := now.Add(time.Second)
	vals, err := s.GetMetricRange("host1", "cpu_percent", from, to)
	if err != nil {
		t.Fatalf("GetMetricRange: %v", err)
	}
	if len(vals) != 1 {
		t.Errorf("expected 1 metric, got %d", len(vals))
	}
	if vals[0].Value != 0.75 {
		t.Errorf("value = %v; want 0.75", vals[0].Value)
	}
}

func TestSaveAndGetAlert(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	alert := map[string]interface{}{
		"id":   "alert1",
		"name": "stress_panic",
	}
	if err := s.SaveAlert("alert1", alert); err != nil {
		t.Fatalf("SaveAlert: %v", err)
	}

	var got map[string]interface{}
	if err := s.GetAlert("alert1", &got); err != nil {
		t.Fatalf("GetAlert: %v", err)
	}
	if got["name"] != "stress_panic" {
		t.Errorf("name = %v; want stress_panic", got["name"])
	}
}

func TestSaveAndGetDashboard(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	type dash struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := s.SaveDashboard("dash1", dash{ID: "dash1", Name: "System Overview"}); err != nil {
		t.Fatalf("SaveDashboard: %v", err)
	}

	var got dash
	if err := s.GetDashboard("dash1", &got); err != nil {
		t.Fatalf("GetDashboard: %v", err)
	}
	if got.Name != "System Overview" {
		t.Errorf("name = %q; want System Overview", got.Name)
	}

	if err := s.DeleteDashboard("dash1"); err != nil {
		t.Fatalf("DeleteDashboard: %v", err)
	}
	if err := s.GetDashboard("dash1", &got); err == nil {
		t.Error("expected error for deleted dashboard")
	}
}

func TestListDashboards(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	for i := 0; i < 3; i++ {
		id := string(rune('a' + i))
		_ = s.SaveDashboard(id, map[string]string{"id": id})
	}

	count := 0
	_ = s.ListDashboards(func(val []byte) error {
		count++
		return nil
	})
	if count != 3 {
		t.Errorf("expected 3 dashboards, got %d", count)
	}
}

func TestSaveAndGetUser(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	type user struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := s.SaveUser("admin", user{Username: "admin", Role: "admin"}); err != nil {
		t.Fatalf("SaveUser: %v", err)
	}

	var got user
	if err := s.GetUser("admin", &got); err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Role != "admin" {
		t.Errorf("role = %q; want admin", got.Role)
	}
}

func TestKPIRangeFilter(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	base := time.Now()
	// Save KPIs at t-2h, t-1h, t
	for i := 0; i < 3; i++ {
		ts := base.Add(-time.Duration(2-i) * time.Hour)
		_ = s.SaveKPI("host1", "stress", float64(i)*0.1, ts)
	}

	// Query only last 90 minutes
	from := base.Add(-90 * time.Minute)
	vals, err := s.GetKPIRange("host1", "stress", from, base.Add(time.Minute))
	if err != nil {
		t.Fatalf("GetKPIRange: %v", err)
	}
	if len(vals) != 2 {
		t.Errorf("expected 2 KPI points in range, got %d", len(vals))
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "ohe-storage-test-*")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return s
}
