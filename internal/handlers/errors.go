package handlers

import (
	"net/http"
)

// notFound renders the friendly 404 page.
func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusNotFound, "404", nil)
}

// serverError logs the real error for us and shows users a friendly 500 page.
// Raw Go errors can leak internals (SQL, file paths), so users never see them.
func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	h.log.Error("server error", "method", r.Method, "path", r.URL.Path, "err", err)

	buf, renderErr := h.tmpl.execute("500", "layout", nil)
	if renderErr != nil {
		// The error page itself failed; fall back to plain text.
		// (Don't call h.render here: it would loop back into serverError.)
		h.log.Error("render 500 page", "err", renderErr)
		http.Error(w, "Something went wrong. Please try again.", http.StatusInternalServerError)
		return
	}
	writeHTML(w, http.StatusInternalServerError, buf)
}
