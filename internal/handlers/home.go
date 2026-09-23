package handlers

import "net/http"

// home shows the placeholder landing page. Slice 1 replaces it with the city picker.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusOK, "home", nil)
}
