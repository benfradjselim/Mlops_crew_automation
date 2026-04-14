package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	oheproc "github.com/benfradjselim/ohe/internal/processor"
	"github.com/benfradjselim/ohe/internal/alerter"
	"github.com/benfradjselim/ohe/internal/analyzer"
	"github.com/benfradjselim/ohe/internal/predictor"
	"github.com/benfradjselim/ohe/internal/storage"
	"github.com/benfradjselim/ohe/pkg/models"
	"github.com/benfradjselim/ohe/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const version = "4.0.0"

// Handlers holds all API dependencies
type Handlers struct {
	store       *storage.Store
	processor   *oheproc.Processor
	analyzer    *analyzer.Analyzer
	predictor   *predictor.Predictor
	alerter     *alerter.Alerter
	hub         *Hub
	jwtSecret   string
	startTime   time.Time
	authEnabled bool
}

// NewHandlers constructs the handler set
func NewHandlers(
	store *storage.Store,
	proc *oheproc.Processor,
	ana *analyzer.Analyzer,
	pred *predictor.Predictor,
	alrt *alerter.Alerter,
	jwtSecret string,
	authEnabled bool,
) *Handlers {
	return &Handlers{
		store:       store,
		processor:   proc,
		analyzer:    ana,
		predictor:   pred,
		alerter:     alrt,
		hub:         NewHub(),
		jwtSecret:   jwtSecret,
		startTime:   time.Now(),
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
	name := mux.Vars(r)["name"]
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
	name := mux.Vars(r)["name"]
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
	name := mux.Vars(r)["name"]
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
	id := mux.Vars(r)["id"]
	al, ok := h.alerter.GetByID(id)
	if !ok {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "alert not found")
		return
	}
	respondSuccess(w, al)
}

// AlertAcknowledgeHandler POST /api/v1/alerts/{id}/acknowledge
func (h *Handlers) AlertAcknowledgeHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.alerter.Acknowledge(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondSuccess(w, map[string]string{"status": "acknowledged"})
}

// AlertSilenceHandler POST /api/v1/alerts/{id}/silence
func (h *Handlers) AlertSilenceHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.alerter.Silence(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondSuccess(w, map[string]string{"status": "silenced"})
}

// AlertDeleteHandler DELETE /api/v1/alerts/{id}
func (h *Handlers) AlertDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
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
	id := mux.Vars(r)["id"]
	var d models.Dashboard
	if err := h.store.GetDashboard(id, &d); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "dashboard not found")
		return
	}
	respondSuccess(w, d)
}

// DashboardUpdateHandler PUT /api/v1/dashboards/{id}
func (h *Handlers) DashboardUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
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
	id := mux.Vars(r)["id"]
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

	// Broadcast live update to WebSocket subscribers
	if msg, err := json.Marshal(map[string]interface{}{
		"type": "kpi_update", "data": snapshot,
	}); err == nil {
		h.hub.Broadcast(msg)
	}

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

// ReloadHandler POST /api/v1/reload — signal config reload (no-op for now, returns ok)
func (h *Handlers) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, map[string]string{"status": "reloaded"})
}

// --- DataSource handlers ---

// DataSourceListHandler GET /api/v1/datasources
func (h *Handlers) DataSourceListHandler(w http.ResponseWriter, r *http.Request) {
	var sources []*models.DataSource
	err := h.store.ListDataSources(func(val []byte) error {
		var ds models.DataSource
		if err := json.Unmarshal(val, &ds); err != nil {
			return nil
		}
		sources = append(sources, &ds)
		return nil
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, sources)
}

// DataSourceCreateHandler POST /api/v1/datasources
func (h *Handlers) DataSourceCreateHandler(w http.ResponseWriter, r *http.Request) {
	var ds models.DataSource
	if err := decodeBody(r, &ds); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	ds.ID = utils.GenerateID(8)
	ds.Enabled = true
	if err := h.store.SaveDataSource(ds.ID, ds); err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": ds, "timestamp": time.Now().UTC(),
	})
}

// DataSourceGetHandler GET /api/v1/datasources/{id}
func (h *Handlers) DataSourceGetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var ds models.DataSource
	if err := h.store.GetDataSource(id, &ds); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "datasource not found")
		return
	}
	respondSuccess(w, ds)
}

// DataSourceUpdateHandler PUT /api/v1/datasources/{id}
func (h *Handlers) DataSourceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var ds models.DataSource
	if err := decodeBody(r, &ds); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	ds.ID = id
	if err := h.store.SaveDataSource(id, ds); err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, ds)
}

// DataSourceDeleteHandler DELETE /api/v1/datasources/{id}
func (h *Handlers) DataSourceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.store.DeleteDataSource(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "datasource not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DataSourceTestHandler POST /api/v1/datasources/{id}/test
func (h *Handlers) DataSourceTestHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var ds models.DataSource
	if err := h.store.GetDataSource(id, &ds); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "datasource not found")
		return
	}
	// Attempt an HTTP GET to the datasource URL
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(ds.URL)
	if err != nil {
		respondSuccess(w, map[string]interface{}{"status": "error", "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	respondSuccess(w, map[string]interface{}{"status": "ok", "http_status": resp.StatusCode})
}

// --- User management handlers ---

// UserListHandler GET /api/v1/auth/users
func (h *Handlers) UserListHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	err := h.store.ListUsers(func(val []byte) error {
		var u models.User
		if err := json.Unmarshal(val, &u); err != nil {
			return nil
		}
		u.Password = "" // never expose hash
		users = append(users, u)
		return nil
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	respondSuccess(w, users)
}

// UserCreateHandler POST /api/v1/auth/users
func (h *Handlers) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := decodeBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "username and password required")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "HASH_ERROR", "could not hash password")
		return
	}
	role := req.Role
	if role == "" {
		role = "viewer"
	}
	user := models.User{
		ID:       utils.GenerateID(8),
		Username: req.Username,
		Password: string(hash),
		Role:     role,
	}
	if err := h.store.SaveUser(req.Username, user); err != nil {
		respondError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}
	user.Password = ""
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": user, "timestamp": time.Now().UTC(),
	})
}

// UserGetHandler GET /api/v1/auth/users/{id}
func (h *Handlers) UserGetHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	var user models.User
	if err := h.store.GetUser(username, &user); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}
	user.Password = ""
	respondSuccess(w, user)
}

// UserDeleteHandler DELETE /api/v1/auth/users/{id}
func (h *Handlers) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	if err := h.store.DeleteUser(username); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LogoutHandler POST /api/v1/auth/logout — stateless JWT, just acknowledge
func (h *Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, map[string]string{"status": "logged out"})
}

// RefreshHandler POST /api/v1/auth/refresh — issue a new token from valid existing one
func (h *Handlers) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := claimsFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no claims in context")
		return
	}
	exp := time.Now().Add(24 * time.Hour)
	newClaims := JWTClaims{
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "TOKEN_ERROR", "could not generate token")
		return
	}
	respondSuccess(w, map[string]interface{}{
		"token":   signed,
		"expires": exp.Unix(),
	})
}

// DashboardExportHandler GET /api/v1/dashboards/{id}/export
func (h *Handlers) DashboardExportHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var d models.Dashboard
	if err := h.store.GetDashboard(id, &d); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "dashboard not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="dashboard-%s.json"`, id))
	_ = json.NewEncoder(w).Encode(d)
}

// DashboardImportHandler POST /api/v1/dashboards/import
func (h *Handlers) DashboardImportHandler(w http.ResponseWriter, r *http.Request) {
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
		"success": true, "data": d, "timestamp": time.Now().UTC(),
	})
}

// QueryHandler POST /api/v1/query — simple metric query by name and time range
func (h *Handlers) QueryHandler(w http.ResponseWriter, r *http.Request) {
	var req models.QueryRequest
	if err := decodeBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	if req.To.IsZero() {
		req.To = time.Now()
	}
	if req.From.IsZero() {
		req.From = req.To.Add(-time.Hour)
	}

	host := r.URL.Query().Get("host")
	if host == "" {
		host = "localhost"
	}

	// req.Query is a metric name for now
	tvs, err := h.store.GetMetricRange(host, req.Query, req.From, req.To)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "QUERY_ERROR", err.Error())
		return
	}

	points := make([]models.DataPoint, 0, len(tvs))
	for _, tv := range tvs {
		points = append(points, models.DataPoint{Timestamp: tv.Timestamp, Value: tv.Value})
	}

	// Downsample if step > 0
	if req.Step > 0 {
		points = oheproc.Downsample(points, time.Duration(req.Step)*time.Second)
	}

	respondSuccess(w, models.QueryResult{Metric: req.Query, Points: points})
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
