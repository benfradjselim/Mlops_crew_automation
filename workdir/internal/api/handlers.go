package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/benfradjselim/ohe/internal/alerter"
	"github.com/benfradjselim/ohe/internal/analyzer"
	"github.com/benfradjselim/ohe/internal/predictor"
	"github.com/benfradjselim/ohe/internal/processor"
	"github.com/benfradjselim/ohe/internal/storage"
	"github.com/benfradjselim/ohe/pkg/models"
	"github.com/benfradjselim/ohe/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const version = "4.0.0"

// Handlers holds all API dependencies
type Handlers struct {
	store     *storage.Store
	processor *processor.Processor
	analyzer  *analyzer.Analyzer
	predictor *predictor.Predictor
	alerter   *alerter.Alerter
	jwtSecret string
	startTime time.Time
	authEnabled bool
}

// NewHandlers constructs the handler set
func NewHandlers(
	store *storage.Store,
	proc *processor.Processor,
	ana *analyzer.Analyzer,
	pred *predictor.Predictor,
	alrt *alerter.Alerter,
	jwtSecret string,
	authEnabled bool,
) *Handlers {
	return &Handlers{
		store:     store,
		processor: proc,
		analyzer:  ana,
		predictor: pred,
		alerter:   alrt,
		jwtSecret: jwtSecret,
		startTime: time.Now(),
		authEnabled: authEnabled,
	}
}

// HealthHandler GET /api/v1/health
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	checks := map[string]string{
		"storage": "ok",
	}
	if !h.store.Healthy() {
		checks["storage"] = "error"
	}

	respondSuccess(w, models.HealthResponse{
		Status:    "ok",
		Version:   version,
		Uptime:    time.Since(h.startTime).Seconds(),
		Checks:    checks,
		Timestamp: time.Now().UTC(),
	})
}

// MetricsListHandler GET /api/v1/metrics
func (h *Handlers) MetricsListHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	// Return latest normalized values for common metrics
	metricNames := []string{
		"cpu_percent", "memory_percent", "disk_percent",
		"net_rx_bps", "net_tx_bps", "load_avg_1", "load_avg_5",
		"load_avg_15", "uptime_seconds", "processes",
	}

	result := make(map[string]interface{})
	for _, name := range metricNames {
		val, ok := h.processor.GetNormalized(host, name)
		if ok {
			result[name] = val
		}
	}

	respondSuccess(w, map[string]interface{}{
		"host":    host,
		"metrics": result,
	})
}

// MetricGetHandler GET /api/v1/metrics/{name}
func (h *Handlers) MetricGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	from, to := parseTimeRange(r)

	values, err := h.store.GetMetricRange(host, name, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}

	respondSuccess(w, map[string]interface{}{
		"host":   host,
		"metric": name,
		"from":   from,
		"to":     to,
		"points": values,
	})
}

// MetricAggregateHandler GET /api/v1/metrics/{name}/aggregate
func (h *Handlers) MetricAggregateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	agg, ok := h.processor.Aggregate(host, name)
	if !ok {
		respondError(w, http.StatusNotFound, "NO_DATA", "no data for metric")
		return
	}
	respondSuccess(w, agg)
}

// KPIListHandler GET /api/v1/kpis
func (h *Handlers) KPIListHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	// Get current normalized values and compute KPIs
	metrics := h.buildMetricsMap(host)
	snapshot := h.analyzer.Update(host, metrics)
	respondSuccess(w, snapshot)
}

// KPIGetHandler GET /api/v1/kpis/{name}
func (h *Handlers) KPIGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	from, to := parseTimeRange(r)

	values, err := h.store.GetKPIRange(host, name, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, map[string]interface{}{
		"host":   host,
		"kpi":    name,
		"from":   from,
		"to":     to,
		"points": values,
	})
}

// PredictHandler GET /api/v1/kpis/{name}/predict or GET /api/v1/predict
func (h *Handlers) PredictHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}
	horizonStr := r.URL.Query().Get("horizon")
	horizon := 120 // default 2 hours
	if horizonStr != "" {
		if v, err := strconv.Atoi(horizonStr); err == nil && v > 0 {
			horizon = v
		}
	}

	preds := h.predictor.PredictAll(host, horizon)
	respondSuccess(w, map[string]interface{}{
		"host":            host,
		"horizon_minutes": horizon,
		"predictions":     preds,
	})
}

// AlertListHandler GET /api/v1/alerts
func (h *Handlers) AlertListHandler(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"
	var alerts []*models.Alert
	if activeOnly {
		alerts = h.alerter.GetActive()
	} else {
		alerts = h.alerter.GetAll()
	}
	respondSuccess(w, alerts)
}

// AlertGetHandler GET /api/v1/alerts/{id}
func (h *Handlers) AlertGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	al, ok := h.alerter.GetByID(id)
	if !ok {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "alert not found")
		return
	}
	respondSuccess(w, al)
}

// AlertAcknowledgeHandler POST /api/v1/alerts/{id}/acknowledge
func (h *Handlers) AlertAcknowledgeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.alerter.Acknowledge(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondSuccess(w, map[string]string{"status": "acknowledged"})
}

// AlertSilenceHandler POST /api/v1/alerts/{id}/silence
func (h *Handlers) AlertSilenceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.alerter.Silence(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondSuccess(w, map[string]string{"status": "silenced"})
}

// AlertDeleteHandler DELETE /api/v1/alerts/{id}
func (h *Handlers) AlertDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.alerter.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DashboardListHandler GET /api/v1/dashboards
func (h *Handlers) DashboardListHandler(w http.ResponseWriter, r *http.Request) {
	var dashboards []*models.Dashboard
	err := h.store.ListDashboards(func(val []byte) error {
		var d models.Dashboard
		if err := json.Unmarshal(val, &d); err != nil {
			return nil
		}
		dashboards = append(dashboards, &d)
		return nil
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, dashboards)
}

// DashboardCreateHandler POST /api/v1/dashboards
func (h *Handlers) DashboardCreateHandler(w http.ResponseWriter, r *http.Request) {
	var d models.Dashboard
	if err := decodeBody(r, &d); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	d.ID = utils.GenerateID(8)
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	if err := h.store.SaveDashboard(d.ID, d); err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"data":      d,
		"timestamp": time.Now().UTC(),
	})
}

// DashboardGetHandler GET /api/v1/dashboards/{id}
func (h *Handlers) DashboardGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d models.Dashboard
	if err := h.store.GetDashboard(id, &d); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "dashboard not found")
		return
	}
	respondSuccess(w, d)
}

// DashboardUpdateHandler PUT /api/v1/dashboards/{id}
func (h *Handlers) DashboardUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d models.Dashboard
	if err := h.store.GetDashboard(id, &d); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "dashboard not found")
		return
	}
	var update models.Dashboard
	if err := decodeBody(r, &update); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	update.ID = id
	update.CreatedAt = d.CreatedAt
	update.UpdatedAt = time.Now()
	if err := h.store.SaveDashboard(id, update); err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, update)
}

// DashboardDeleteHandler DELETE /api/v1/dashboards/{id}
func (h *Handlers) DashboardDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.DeleteDashboard(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "dashboard not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LoginHandler POST /api/v1/auth/login
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decodeBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	var user models.User
	if err := h.store.GetUser(req.Username, &user); err != nil {
		respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
		return
	}

	exp := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "TOKEN_ERROR", "could not generate token")
		return
	}

	respondSuccess(w, models.LoginResponse{
		Token:   signed,
		Expires: exp.Unix(),
		User:    user,
	})
}

// IngestHandler POST /api/v1/ingest — receives metrics from remote agents
func (h *Handlers) IngestHandler(w http.ResponseWriter, r *http.Request) {
	var batch models.MetricBatch
	if err := decodeBody(r, &batch); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	h.processor.Ingest(batch.Metrics)

	// Store metrics in Badger
	for _, m := range batch.Metrics {
		_ = h.store.SaveMetric(m.Host, m.Name, m.Value, m.Timestamp)
	}

	// Build metrics map and run KPI analysis
	metrics := h.buildMetricsMap(batch.Host)
	snapshot := h.analyzer.Update(batch.Host, metrics)

	// Store KPIs
	_ = h.store.SaveKPI(batch.Host, "stress", snapshot.Stress.Value, snapshot.Timestamp)
	_ = h.store.SaveKPI(batch.Host, "fatigue", snapshot.Fatigue.Value, snapshot.Timestamp)
	_ = h.store.SaveKPI(batch.Host, "mood", snapshot.Mood.Value, snapshot.Timestamp)
	_ = h.store.SaveKPI(batch.Host, "pressure", snapshot.Pressure.Value, snapshot.Timestamp)
	_ = h.store.SaveKPI(batch.Host, "humidity", snapshot.Humidity.Value, snapshot.Timestamp)
	_ = h.store.SaveKPI(batch.Host, "contagion", snapshot.Contagion.Value, snapshot.Timestamp)

	// Feed predictor
	now := time.Now()
	for _, m := range batch.Metrics {
		h.predictor.Feed(m.Host, m.Name, m.Value, now)
	}
	h.predictor.Feed(batch.Host, "stress", snapshot.Stress.Value, now)
	h.predictor.Feed(batch.Host, "fatigue", snapshot.Fatigue.Value, now)

	// Evaluate alerts
	kpiMap := map[string]float64{
		"stress":    snapshot.Stress.Value,
		"fatigue":   snapshot.Fatigue.Value,
		"mood":      snapshot.Mood.Value,
		"pressure":  snapshot.Pressure.Value,
		"humidity":  snapshot.Humidity.Value,
		"contagion": snapshot.Contagion.Value,
	}
	h.alerter.Evaluate(batch.Host, kpiMap)

	respondSuccess(w, map[string]interface{}{
		"accepted": len(batch.Metrics),
		"kpis":     snapshot,
	})
}

// ConfigHandler GET /api/v1/config
func (h *Handlers) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, map[string]interface{}{
		"version":      version,
		"auth_enabled": h.authEnabled,
	})
}

// --- Helpers ---

func (h *Handlers) buildMetricsMap(host string) map[string]float64 {
	names := []string{
		"cpu_percent", "memory_percent", "disk_percent",
		"load_avg_1", "error_rate", "timeout_rate",
		"request_rate", "uptime_seconds",
	}
	m := make(map[string]float64, len(names))
	for _, name := range names {
		if v, ok := h.processor.GetNormalized(host, name); ok {
			m[name] = v
		}
	}
	return m
}

func parseTimeRange(r *http.Request) (from, to time.Time) {
	now := time.Now()
	to = now

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if from.IsZero() {
		from = now.Add(-time.Hour)
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}
	return from, to
}

func decodeBody(r *http.Request, dest interface{}) error {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10MB limit
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}
