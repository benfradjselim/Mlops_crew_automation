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

	api := r.PathPrefix("/api/v1").Subrouter()

	// System
	api.HandleFunc("/health", h.HealthHandler).Methods(http.MethodGet)
	api.HandleFunc("/config", h.ConfigHandler).Methods(http.MethodGet)

	// Auth
	api.HandleFunc("/auth/login", h.LoginHandler).Methods(http.MethodPost)

	// Metrics
	api.HandleFunc("/metrics", h.MetricsListHandler).Methods(http.MethodGet)
	api.HandleFunc("/metrics/{name}", h.MetricGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/metrics/{name}/aggregate", h.MetricAggregateHandler).Methods(http.MethodGet)

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
	api.HandleFunc("/dashboards/{id}", h.DashboardGetHandler).Methods(http.MethodGet)
	api.HandleFunc("/dashboards/{id}", h.DashboardUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/dashboards/{id}", h.DashboardDeleteHandler).Methods(http.MethodDelete)

	// Ingest (agent → central push endpoint)
	api.HandleFunc("/ingest", h.IngestHandler).Methods(http.MethodPost)

	// Embedded UI — serve static files if present
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return r
}
