package models

import "time"

// Metric represents a single time-series data point
type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels,omitempty"`
	Host      string            `json:"host"`
}

// MetricBatch is a collection of metrics sent by an agent
type MetricBatch struct {
	AgentID   string    `json:"agent_id"`
	Host      string    `json:"host"`
	Metrics   []Metric  `json:"metrics"`
	Timestamp time.Time `json:"timestamp"`
}

// KPI represents a computed composite KPI
type KPI struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	State     string    `json:"state"`
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
}

// KPISnapshot holds all current KPIs for a host
type KPISnapshot struct {
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
	Stress    KPI       `json:"stress"`
	Fatigue   KPI       `json:"fatigue"`
	Mood      KPI       `json:"mood"`
	Pressure  KPI       `json:"pressure"`
	Humidity  KPI       `json:"humidity"`
	Contagion KPI       `json:"contagion"`
}

// Prediction is a forecasted value for a metric/KPI
type Prediction struct {
	Target    string    `json:"target"`
	Current   float64   `json:"current"`
	Predicted float64   `json:"predicted"`
	Horizon   int       `json:"horizon_minutes"`
	Trend     string    `json:"trend"` // "rising", "stable", "falling"
	Timestamp time.Time `json:"timestamp"`
}

// Alert severity levels
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
	SeverityEmergency = "emergency"
)

// Alert status
const (
	StatusActive       = "active"
	StatusAcknowledged = "acknowledged"
	StatusSilenced     = "silenced"
	StatusResolved     = "resolved"
)

// Alert represents a triggered observability alert
type Alert struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Host        string    `json:"host"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Prediction  string    `json:"prediction,omitempty"` // e.g. "Storm in 2 hours"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// User for auth
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // bcrypt hash — omitempty ensures it is never serialised when cleared
	Role     string `json:"role"` // admin, viewer, operator
}

// Dashboard configuration
type Dashboard struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Widgets     []Widget  `json:"widgets"`
	Refresh     int       `json:"refresh_seconds"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Widget types
const (
	WidgetTypeTimeseries = "timeseries"
	WidgetTypeGauge      = "gauge"
	WidgetTypeKPI        = "kpi"
	WidgetTypeStat       = "stat"
	WidgetTypeAlerts     = "alerts"
)

// Widget is a dashboard panel
type Widget struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Metric      string            `json:"metric,omitempty"`
	KPI         string            `json:"kpi,omitempty"`
	Aggregation string            `json:"aggregation,omitempty"` // avg, min, max, p95, p99
	From        string            `json:"from,omitempty"`        // relative: -1h, -24h
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	Options     map[string]string `json:"options,omitempty"`
}

// DataSource represents an external data source
type DataSource struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"` // prometheus, loki, custom
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers,omitempty"`
	Enabled  bool              `json:"enabled"`
}

// SystemMetrics holds all raw collected system metrics
type SystemMetrics struct {
	Host         string    `json:"host"`
	Timestamp    time.Time `json:"timestamp"`
	CPUPercent   float64   `json:"cpu_percent"`
	MemoryPercent float64  `json:"memory_percent"`
	MemoryUsedMB float64   `json:"memory_used_mb"`
	MemoryTotalMB float64  `json:"memory_total_mb"`
	DiskPercent  float64   `json:"disk_percent"`
	DiskUsedGB   float64   `json:"disk_used_gb"`
	DiskTotalGB  float64   `json:"disk_total_gb"`
	NetRxBps     float64   `json:"net_rx_bps"`
	NetTxBps     float64   `json:"net_tx_bps"`
	LoadAvg1     float64   `json:"load_avg_1"`
	LoadAvg5     float64   `json:"load_avg_5"`
	LoadAvg15    float64   `json:"load_avg_15"`
	Processes    int       `json:"processes"`
	Uptime       float64   `json:"uptime_seconds"`
}

// APIResponse wraps all API responses
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError contains error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// QueryRequest for QQL queries
type QueryRequest struct {
	Query string    `json:"query"`
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Step  int       `json:"step_seconds"`
}

// QueryResult holds query results
type QueryResult struct {
	Metric string      `json:"metric"`
	Points []DataPoint `json:"points"`
}

// DataPoint is a single time-value pair
type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// LoginRequest for auth endpoint
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains JWT token
type LoginResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	User    User   `json:"user"`
}

// HealthResponse for health check
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Host      string            `json:"host"`
	Uptime    float64           `json:"uptime_seconds"`
	Checks    map[string]string `json:"checks"`
	Timestamp time.Time         `json:"timestamp"`
}
