package httpapi

import "net/http"

// triggerSync kicks off a full catalog re-sync in the background and returns
// immediately, since a full sync takes real wall-clock minutes.
func (h *handlers) triggerSync(w http.ResponseWriter, r *http.Request) {
	if h.deps.SyncTrigger != nil {
		go h.deps.SyncTrigger()
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"message": "sync started"})
}
