package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter builds and returns the HTTP router with all routes registered
func NewRouter(h *Handlers, jwtSecret string, authEnabled bool) http.Handler {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(LoggingMiddleware)
	r.Use(CORSMiddleware)
	r.Use(mux.MiddlewareFunc(AuthMiddleware(jwtSecret, authEnabled)))
	r.Use(mux.MiddlewareFunc(RateLimitLogin))

	api := r.PathPrefix("/api/v1").Subrouter()

	// System
	api.HandleFunc("/health", h.HealthHandler).Methods(http.MethodGet)
	api.HandleFunc("/config", h.ConfigHandler).Methods(http.MethodGet)

	// Auth
	api.HandleFunc("/auth/login", h.LoginHandler).Methods(http.MethodPost)
	api.HandleFunc("/auth/logout", h.LogoutHandler).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", h.RefreshHandler).Methods(http.MethodPost)
	adminOnly := RequireRole("admin")
	api.Handle("/auth/users", adminOnly(http.HandlerFunc(h.UserListHandler))).Methods(http.MethodGet)
	api.Handle("/auth/users", adminOnly(http.HandlerFunc(h.UserCreateHandler))).Methods(http.MethodPost)
	api.Handle("/auth/users/{id}", adminOnly(http.HandlerFunc(h.UserGetHandler))).Methods(http.MethodGet)
	api.Handle("/auth/users/{id}", adminOnly(http.HandlerFunc(h.UserDeleteHandler))).Methods(http.MethodDelete)
	api.Handle("/reload", adminOnly(http.HandlerFunc(h.ReloadHandler))).Methods(http.MethodPost)

	// Metrics
	api.HandleFunc("/metrics", h.MetricsListHandler).Methods(http.MethodGet)
	api.HandleFunc("/metrics/{name}", h.MetricGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/metrics/{name}/aggregate", h.MetricAggregateHandler).Methods(http.MethodGet)

	// Query (QQL)
	api.HandleFunc("/query", h.QueryHandler).Methods(http.MethodPost)

	// KPIs
	api.HandleFunc("/kpis", h.KPIListHandler).Methods(http.MethodGet)
	api.HandleFunc("/kpis/{name}", h.KPIGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/kpis/{name}/predict", h.PredictHandler).Methods(http.MethodGet)
	api.HandleFunc("/predict", h.PredictHandler).Methods(http.MethodGet)

	// Alerts
	api.HandleFunc("/alerts", h.AlertListHandler).Methods(http.MethodGet)
	api.HandleFunc("/alerts/{id}", h.AlertGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/alerts/{id}", h.AlertDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/alerts/{id}/acknowledge", h.AlertAcknowledgeHandler).Methods(http.MethodPost)
	api.HandleFunc("/alerts/{id}/silence", h.AlertSilenceHandler).Methods(http.MethodPost)

	// Dashboards
	api.HandleFunc("/dashboards", h.DashboardListHandler).Methods(http.MethodGet)
	api.HandleFunc("/dashboards", h.DashboardCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/dashboards/import", h.DashboardImportHandler).Methods(http.MethodPost)
	api.HandleFunc("/dashboards/{id}", h.DashboardGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/dashboards/{id}", h.DashboardUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/dashboards/{id}", h.DashboardDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/dashboards/{id}/export", h.DashboardExportHandler).Methods(http.MethodGet)

	// DataSources
	api.HandleFunc("/datasources", h.DataSourceListHandler).Methods(http.MethodGet)
	api.HandleFunc("/datasources", h.DataSourceCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/datasources/{id}", h.DataSourceGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/datasources/{id}", h.DataSourceUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/datasources/{id}", h.DataSourceDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/datasources/{id}/test", h.DataSourceTestHandler).Methods(http.MethodPost)

	// Ingest (agent → central push endpoint)
	api.HandleFunc("/ingest", h.IngestHandler).Methods(http.MethodPost)

	// WebSocket streaming
	api.HandleFunc("/ws", h.WebSocketHandler)

	// Embedded UI
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return r
}
