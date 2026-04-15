package receiver

import (
	"log"

	"github.com/benfradjselim/ohe/internal/analyzer"
	"github.com/benfradjselim/ohe/internal/storage"
	"github.com/benfradjselim/ohe/pkg/models"
)

// Bus wires ingested data from all receivers into the OHE pipeline.
// It implements MetricSink, SpanSink, and LogSink.
type Bus struct {
	store    *storage.Store
	topology *analyzer.TopologyAnalyzer
}

// NewBus creates a receiver bus wired to the given store and topology analyzer
func NewBus(store *storage.Store, topology *analyzer.TopologyAnalyzer) *Bus {
	return &Bus{store: store, topology: topology}
}

// IngestMetric satisfies MetricSink — forwards metric to the store
func (b *Bus) IngestMetric(m models.Metric) {
	if err := b.store.SaveMetric(m.Host, m.Name, m.Value, m.Timestamp); err != nil {
		log.Printf("[bus] metric %s/%s: %v", m.Host, m.Name, err)
	}
}

// IngestSpan satisfies SpanSink — forwards span to topology analyzer
func (b *Bus) IngestSpan(s models.Span) {
	if b.topology != nil {
		b.topology.IngestSpan(s)
	}
	if err := b.store.SaveSpan(s, s.TraceID, s.SpanID); err != nil {
		log.Printf("[bus] span %s/%s: %v", s.TraceID, s.SpanID, err)
	}
}

// IngestLog satisfies LogSink — stores log entry
func (b *Bus) IngestLog(e models.LogEntry) {
	service := e.Service
	if service == "" {
		service = e.Host
	}
	if service == "" {
		service = "unknown"
	}
	if err := b.store.SaveLog(service, e, e.Timestamp); err != nil {
		log.Printf("[bus] log %s: %v", service, err)
	}
}
