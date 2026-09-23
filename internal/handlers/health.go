package handlers

import (
	"context"
	"net/http"
	"time"
)

// health reports whether the app can reach the database. Monitoring tools
// and Docker call this, so it returns plain text rather than an HTML page.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	// Don't let a hung database make the health check hang too.
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err := h.store.Ping(ctx); err != nil {
		h.log.Error("health check failed", "err", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("database unavailable"))
		return
	}
	_, _ = w.Write([]byte("ok"))
}
