package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// Templates holds one parsed template set per page.
//
// Why one set per page? Every page defines blocks with the same names
// ("title", "content"). If all pages were parsed into a single set, the last
// page parsed would overwrite the others' blocks. Giving each page its own
// copy of layout + partials avoids that.
type Templates struct {
	pages map[string]*template.Template
}

// LoadTemplates parses templates/layout.html, any shared partials in
// templates/partials/, and each page in templates/pages/.
// A page named "home" comes from templates/pages/home.html.
func LoadTemplates(fsys fs.FS) (*Templates, error) {
	shared := []string{"templates/layout.html"}
	partials, err := fs.Glob(fsys, "templates/partials/*.html")
	if err != nil {
		return nil, err
	}
	shared = append(shared, partials...)

	pageFiles, err := fs.Glob(fsys, "templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	t := &Templates{pages: make(map[string]*template.Template)}
	for _, file := range pageFiles {
		name := strings.TrimSuffix(path.Base(file), ".html")
		files := append(append([]string{}, shared...), file)
		tmpl, err := template.ParseFS(fsys, files...)
		if err != nil {
			return nil, fmt.Errorf("parse page %q: %w", name, err)
		}
		t.pages[name] = tmpl
	}
	return t, nil
}

// execute renders a named template from a page's set into a buffer.
// Rendering into a buffer first (instead of straight to the response) means a
// template error can't leave the user with half a page.
func (t *Templates) execute(page, tmplName string, data any) (*bytes.Buffer, error) {
	tmpl, ok := t.pages[page]
	if !ok {
		return nil, fmt.Errorf("unknown page %q", page)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, data); err != nil {
		return nil, fmt.Errorf("render %s/%s: %w", page, tmplName, err)
	}
	return &buf, nil
}

// render writes a full page (layout + page content) with the given status.
func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data any) {
	buf, err := h.tmpl.execute(page, "layout", data)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	writeHTML(w, status, buf)
}

// renderPartial writes just one named template (no layout), for HTMX to swap
// into the page. page picks which template set to use; partials are in all sets.
func (h *Handler) renderPartial(w http.ResponseWriter, r *http.Request, status int, page, name string, data any) {
	buf, err := h.tmpl.execute(page, name, data)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	writeHTML(w, status, buf)
}

func writeHTML(w http.ResponseWriter, status int, buf *bytes.Buffer) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w) // nothing useful to do if the client hung up
}
