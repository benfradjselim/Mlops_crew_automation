package predictor

import (
	"math"
	"testing"
	"time"
)

func TestILRBasic(t *testing.T) {
	m := NewILR()
	if m.IsTrained() {
		t.Error("fresh model should not be trained")
	}

	// Feed y = 2x + 1
	for i := 0; i < 20; i++ {
		x := float64(i)
		m.Update(x, 2*x+1)
	}

	if !m.IsTrained() {
		t.Error("model should be trained after 20 updates")
	}

	// Alpha ≈ 2.0, Beta ≈ 1.0
	if math.Abs(m.Alpha-2.0) > 0.01 {
		t.Errorf("Alpha = %v; want ~2.0", m.Alpha)
	}
	if math.Abs(m.Beta-1.0) > 0.01 {
		t.Errorf("Beta = %v; want ~1.0", m.Beta)
	}

	// Predict x=10 → 21
	pred := m.Predict(10)
	if math.Abs(pred-21.0) > 0.1 {
		t.Errorf("Predict(10) = %v; want ~21", pred)
	}
}

func TestILRTrend(t *testing.T) {
	m := NewILR()
	// Rising trend
	for i := 0; i < 10; i++ {
		m.Update(float64(i), float64(i)*2)
	}
	if m.Trend() != "rising" {
		t.Errorf("Trend() = %q; want rising", m.Trend())
	}

	m.Reset()
	// Falling trend
	for i := 0; i < 10; i++ {
		m.Update(float64(i), float64(10-i))
	}
	if m.Trend() != "falling" {
		t.Errorf("Trend() = %q; want falling", m.Trend())
	}
}

func TestBatchILR(t *testing.T) {
	b := NewBatchILR(5)
	for i := 0; i < 25; i++ {
		b.Update(float64(i), float64(i)*3)
	}
	// After 5 full batches, model should predict y=3x
	pred := b.Predict(10)
	if math.Abs(pred-30.0) > 1.0 {
		t.Errorf("BatchILR Predict(10) = %v; want ~30", pred)
	}
}

func TestPredictorFeedAndPredict(t *testing.T) {
	p := NewPredictor()
	now := time.Now()

	// Feed 30 points for a rising metric
	for i := 0; i < 30; i++ {
		ts := now.Add(time.Duration(i) * 15 * time.Second)
		p.Feed("host1", "cpu_percent", float64(i)*2, ts)
	}

	pred, ok := p.Predict("host1", "cpu_percent", 60)
	if !ok {
		t.Fatal("Predict returned not-ok for known metric")
	}
	if pred.Trend != "rising" {
		t.Errorf("expected rising trend, got %q", pred.Trend)
	}
	if pred.Predicted <= pred.Current {
		t.Errorf("predicted %v should be > current %v for rising trend", pred.Predicted, pred.Current)
	}
}

func TestPredictorPredictAll(t *testing.T) {
	p := NewPredictor()
	now := time.Now()

	for i := 0; i < 25; i++ {
		ts := now.Add(time.Duration(i) * 15 * time.Second)
		p.Feed("host2", "cpu_percent", float64(i), ts)
		p.Feed("host2", "memory_percent", float64(50+i), ts)
	}

	preds := p.PredictAll("host2", 30)
	if len(preds) == 0 {
		t.Error("PredictAll returned empty slice")
	}
}
