package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Master-004/glutenscout/internal/models"
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

// funcs are small helpers templates can call, e.g. {{initial .AuthorName}}.
var funcs = template.FuncMap{
	"initial":       initial,
	"stars":         stars,
	"ratingText":    ratingText,
	"reviewerLabel": reviewerLabel,
	"date":          func(t time.Time) string { return t.Format("Jan 2, 2006") },
	"asPhrase":      asPhrase,
	"telDigits":     telDigits,
	"mealImage":     mealImage,
}

// useIllustratedMeals switches meal images to the drawings in
// web/static/images/meals/illustrated/. Set it to false to go back to the
// original simple images; nothing else (including the database) changes.
const useIllustratedMeals = true

const (
	mealImageDir       = "/static/images/meals/"
	illustratedMealDir = "/static/images/meals/illustrated/"
)

// mealImage returns the image URL to show for a meal's photo_url.
// An empty photo_url gets the placeholder. With useIllustratedMeals, our own
// meal images are swapped for the illustrated versions (same file names);
// any other URL (e.g. a real photo, later) is used as is.
func mealImage(photoURL string) string {
	dir := mealImageDir
	if useIllustratedMeals {
		dir = illustratedMealDir
	}
	if photoURL == "" {
		return dir + "placeholder.svg"
	}
	if name, ok := strings.CutPrefix(photoURL, mealImageDir); ok && !strings.Contains(name, "/") {
		return dir + name
	}
	return photoURL
}

// asPhrase turns a badge label into words that fit mid-sentence:
// "Shared oven?" -> "shared oven". Acronyms keep their case ("GF hoisin").
func asPhrase(label string) string {
	label = strings.TrimSuffix(strings.TrimSpace(label), "?")
	r := []rune(label)
	if len(r) >= 2 && unicode.IsUpper(r[0]) && !unicode.IsUpper(r[1]) {
		r[0] = unicode.ToLower(r[0])
	}
	return string(r)
}

// telDigits keeps only the digits of a phone number, for a tel: link:
// "(801) 555-0101" -> "8015550101".
func telDigits(phone string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1 // drop it
	}, phone)
}

// stars turns a rating into five star states ("full", "half" or "empty") for
// the StarRating partial. The design system shows half-star steps only, so
// the rating is rounded to the nearest 0.5 first (4.3 -> 4.5, 4.2 -> 4).
// It takes float64 averages and int ratings from single reviews.
func stars(rating any) []string {
	var r float64
	switch v := rating.(type) {
	case float64:
		r = v
	case int:
		r = float64(v)
	}
	halves := int(math.Round(r * 2))
	out := make([]string, 5)
	for i := range out {
		switch {
		case halves >= 2*(i+1):
			out[i] = "full"
		case halves == 2*i+1:
			out[i] = "half"
		default:
			out[i] = "empty"
		}
	}
	return out
}

// ratingText formats a rating to one decimal place, matching the stars
// (rounded to the nearest half): 4.33 -> "4.5", 4 -> "4.0".
func ratingText(rating float64) string {
	return strconv.FormatFloat(math.Round(rating*2)/2, 'f', 1, 64)
}

// reviewerLabel is the short tag shown next to a reviewer's name.
func reviewerLabel(t string) string {
	switch t {
	case models.ReviewerCeliac:
		return "Celiac"
	case models.ReviewerStrict:
		return "Strict GF"
	case models.ReviewerSensitive:
		return "Gluten-sensitive"
	case models.ReviewerRestaurant:
		return "Restaurant"
	}
	return "" // reviews from before reviewer types existed: no tag
}

// initial returns the first letter of name, uppercased, for avatar monograms.
// It works on runes (not bytes) so names like "Émile" don't get cut mid-character.
func initial(name string) string {
	for _, r := range strings.TrimSpace(name) {
		return strings.ToUpper(string(r))
	}
	return "?"
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
		// Funcs must be added before parsing, so the templates can call them.
		tmpl, err := template.New(path.Base(files[0])).Funcs(funcs).ParseFS(fsys, files...)
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

// layoutData is what the shared layout receives. Page is the page's own data,
// which the layout passes on to the page's "title" and "content" blocks.
// Wrapping it means the layout can know the theme without every page's data
// struct needing a Theme field.
type layoutData struct {
	Theme string // "light", "dark", or "system"
	Page  any
}

// render writes a full page (layout + page content) with the given status.
func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data any) {
	buf, err := h.tmpl.execute(page, "layout", layoutData{Theme: themeFromRequest(r), Page: data})
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
