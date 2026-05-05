package alerter

import (
	"testing"
	"time"
)

func TestIsSuppressed_insideWindow(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()
	a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "default/host/myhost",
		From:        now.Add(-5 * time.Minute),
		Until:       now.Add(5 * time.Minute),
		Reason:      "deploy",
	})

	if !a.isSuppressed("default/host/myhost", now) {
		t.Error("expected isSuppressed=true inside window")
	}
}

func TestIsSuppressed_outsideWindow(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()
	// Window already expired
	a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "default/host/myhost",
		From:        now.Add(-10 * time.Minute),
		Until:       now.Add(-1 * time.Minute),
	})

	if a.isSuppressed("default/host/myhost", now) {
		t.Error("expected isSuppressed=false after window expired")
	}
}

func TestIsSuppressed_wildcardSuppressesAll(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()
	a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "*",
		From:        now.Add(-1 * time.Minute),
		Until:       now.Add(10 * time.Minute),
	})

	if !a.isSuppressed("some/arbitrary/host", now) {
		t.Error("expected wildcard window to suppress all workloads")
	}
}

func TestIsSuppressed_differentWorkloadNotSuppressed(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()
	a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "default/host/hostA",
		From:        now.Add(-1 * time.Minute),
		Until:       now.Add(10 * time.Minute),
	})

	if a.isSuppressed("default/host/hostB", now) {
		t.Error("expected different workload key to not be suppressed")
	}
}

func TestEvaluate_suppressedDuringWindow(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()
	a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "default/host/host5",
		From:        now.Add(-1 * time.Minute),
		Until:       now.Add(10 * time.Minute),
	})

	// Drain channel
	for len(a.ch) > 0 {
		<-a.ch
	}

	a.Evaluate("host5", map[string]float64{"stress": 0.99})

	if len(a.ch) > 0 {
		t.Error("expected no alerts during maintenance window")
	}
}

func TestAddRemoveListMaintenanceWindows(t *testing.T) {
	a := NewAlerter(100)
	now := time.Now()

	id := a.AddMaintenanceWindow(MaintenanceWindow{
		WorkloadKey: "ns/svc",
		From:        now.Add(-1 * time.Minute),
		Until:       now.Add(30 * time.Minute),
		Reason:      "rolling restart",
	})
	if id == "" {
		t.Fatal("expected non-empty ID")
	}

	windows := a.ListMaintenanceWindows()
	if len(windows) != 1 {
		t.Fatalf("expected 1 active window, got %d", len(windows))
	}
	if windows[0].Reason != "rolling restart" {
		t.Errorf("wrong reason: %q", windows[0].Reason)
	}

	a.RemoveMaintenanceWindow(id)
	if len(a.ListMaintenanceWindows()) != 0 {
		t.Error("expected 0 windows after remove")
	}
}
