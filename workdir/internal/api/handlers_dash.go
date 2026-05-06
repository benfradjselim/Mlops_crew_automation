package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// handleHistory returns time-series history for all workloads or a specific one.
func (h *Handlers) handleHistory(w http.ResponseWriter, r *http.Request) {
	if h.historyMgr == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}
	vars := mux.Vars(r)
	if key := vars["workload"]; key != "" {
		writeJSON(w, http.StatusOK, h.historyMgr.Get(key))
		return
	}
	writeJSON(w, http.StatusOK, h.historyMgr.All())
}

// handleEvents returns recent events from the event bus.
func (h *Handlers) handleEvents(w http.ResponseWriter, r *http.Request) {
	if h.eventBus == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	n := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			if parsed > 1000 {
				parsed = 1000
			}
			n = parsed
		}
	}
	writeJSON(w, http.StatusOK, h.eventBus.Recent(n))
}
