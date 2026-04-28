package explain

import (
	"fmt"
	"sync"
	"time"
)

type MetricContribution struct {
	Metric   string
	Pipeline string // "metric"|"log"|"trace"
	Weight   float64
	RValue   float64
}

type ExplainResponse struct {
	RuptureID     string
	Host          string
	R             float64
	Confidence    float64
	Timestamp     time.Time
	Contributions []MetricContribution
	FirstPipeline string // pipeline that fired first ("metric"|"log"|"trace")
}

type FormulaAuditResponse struct {
	RuptureID    string
	AlphaBurst   float64
	AlphaStable  float64
	RuptureIndex float64
	TTFSeconds   float64
	Confidence   float64
	FusedR       float64
	MetricR      float64
	LogR         float64
	TraceR       float64
}

type PipelineDebugResponse struct {
	RuptureID string
	MetricR   float64
	LogR      float64
	TraceR    float64
	FusedR    float64
	Timestamp time.Time
}

type Explainer interface {
	Explain(ruptureID string) (*ExplainResponse, error)
	FormulaAudit(ruptureID string) (*FormulaAuditResponse, error)
	PipelineDebug(ruptureID string) (*PipelineDebugResponse, error)
}

type RuptureRecord struct {
	ID          string
	Host        string
	R           float64
	Confidence  float64
	Timestamp   time.Time
	AlphaBurst  float64
	AlphaStable float64
	TTFSeconds  float64
	MetricR     float64
	LogR        float64
	TraceR      float64
	FusedR      float64
	Metrics     []MetricContribution
}

type Engine struct {
	records sync.Map
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Record(rec RuptureRecord) {
	e.records.Store(rec.ID, rec)
}

func (e *Engine) Explain(id string) (*ExplainResponse, error) {
	val, ok := e.records.Load(id)
	if !ok {
		return nil, fmt.Errorf("explain: rupture %s not found", id)
	}
	rec := val.(RuptureRecord)

	// Determine FirstPipeline
	first := "metric"
	max := rec.MetricR
	if rec.LogR > max {
		first = "log"
		max = rec.LogR
	}
	if rec.TraceR > max {
		first = "trace"
	}

	// Normalize contributions
	var sum float64
	for _, m := range rec.Metrics {
		sum += m.Weight
	}
	var normalized []MetricContribution
	for _, m := range rec.Metrics {
		c := m
		if sum > 0 {
			c.Weight = c.Weight / sum
		}
		normalized = append(normalized, c)
	}

	return &ExplainResponse{
		RuptureID:     rec.ID,
		Host:          rec.Host,
		R:             rec.R,
		Confidence:    rec.Confidence,
		Timestamp:     rec.Timestamp,
		Contributions: normalized,
		FirstPipeline: first,
	}, nil
}

func (e *Engine) FormulaAudit(id string) (*FormulaAuditResponse, error) {
	val, ok := e.records.Load(id)
	if !ok {
		return nil, fmt.Errorf("explain: rupture %s not found", id)
	}
	rec := val.(RuptureRecord)
	return &FormulaAuditResponse{
		RuptureID:    rec.ID,
		AlphaBurst:   rec.AlphaBurst,
		AlphaStable:  rec.AlphaStable,
		RuptureIndex: rec.R,
		TTFSeconds:   rec.TTFSeconds,
		Confidence:   rec.Confidence,
		FusedR:       rec.FusedR,
		MetricR:      rec.MetricR,
		LogR:         rec.LogR,
		TraceR:       rec.TraceR,
	}, nil
}

func (e *Engine) PipelineDebug(id string) (*PipelineDebugResponse, error) {
	val, ok := e.records.Load(id)
	if !ok {
		return nil, fmt.Errorf("explain: rupture %s not found", id)
	}
	rec := val.(RuptureRecord)
	return &PipelineDebugResponse{
		RuptureID: rec.ID,
		MetricR:   rec.MetricR,
		LogR:      rec.LogR,
		TraceR:    rec.TraceR,
		FusedR:    rec.FusedR,
		Timestamp: rec.Timestamp,
	}, nil
}

// NarrativeExplain returns a human-readable explanation of a rupture event.
// It is a structured template filled from the rupture record — no LLM required.
func (e *Engine) NarrativeExplain(id string) (string, error) {
	val, ok := e.records.Load(id)
	if !ok {
		return "", fmt.Errorf("explain: rupture %s not found", id)
	}
	rec := val.(RuptureRecord)

	// Determine primary pipeline
	primaryPipeline := "metric"
	primaryR := rec.MetricR
	if rec.LogR > primaryR {
		primaryPipeline = "log"
		primaryR = rec.LogR
	}
	if rec.TraceR > primaryR {
		primaryPipeline = "trace"
		primaryR = rec.TraceR
	}

	// Severity label
	severity := "warning"
	if rec.R >= 5.0 {
		severity = "critical"
	} else if rec.R >= 3.0 {
		severity = "elevated"
	}

	// Find the top contributing metric
	topMetric := "unknown"
	topWeight := 0.0
	for _, m := range rec.Metrics {
		if m.Weight > topWeight {
			topMetric = m.Metric
			topWeight = m.Weight
		}
	}

	// Build TTF description
	ttfDesc := ""
	if rec.TTFSeconds > 0 {
		mins := int(rec.TTFSeconds / 60)
		if mins < 1 {
			ttfDesc = fmt.Sprintf(" TTF was %ds.", int(rec.TTFSeconds))
		} else {
			ttfDesc = fmt.Sprintf(" TTF was %d minutes.", mins)
		}
	}

	// Contagion note
	contagionNote := ""
	if rec.LogR > 1.0 {
		contagionNote = " Log burst signals indicate contagion may have spread from a dependency."
	}
	if rec.TraceR > 1.0 {
		contagionNote = " Trace error propagation detected — check service dependency graph."
	}

	narrative := fmt.Sprintf(
		"[%s] Rupture %s on %s — R=%.2f (%s). "+
			"Primary signal: %s pipeline (R=%.2f). "+
			"Top contributing factor: %s (weight=%.0f%%)."+
			"%s%s",
		rec.Timestamp.UTC().Format("2006-01-02 15:04:05 UTC"),
		rec.ID,
		rec.Host,
		rec.R, severity,
		primaryPipeline, primaryR,
		topMetric, topWeight*100,
		ttfDesc, contagionNote,
	)
	return narrative, nil
}
