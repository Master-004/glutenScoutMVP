package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/Master-004/glutenscout/internal/store"
)

// maxReviewFormBytes caps the size of a review POST body. 1000 characters can
// be up to 4000 bytes in UTF-8, so 16 KB is plenty and stops huge uploads.
const maxReviewFormBytes = 16 << 10

// upvoteButton is the data for the upvote_button partial.
type upvoteButton struct {
	MealID  int64
	Count   int
	Upvoted bool // highlights the button when this visitor has upvoted
}

// mealCard is one meal on the restaurant page, with its community stats.
type mealCard struct {
	Meal        models.Meal
	Upvote      upvoteButton
	ReviewCount int
}

// reviewForm holds what the visitor typed plus any error messages, so a
// failed submit can re-show the form with their text still in it.
type reviewForm struct {
	MealID    int64
	Author    string
	Body      string
	AuthorErr string
	BodyErr   string
	FormErr   string // problems not tied to one field (e.g. rate limit)
	Posted    bool   // true right after a successful post: shows "Thanks!"
}

// reviewsSection is the data for the reviews_section partial: the form plus
// the list. HTMX re-renders this whole section after each review POST.
type reviewsSection struct {
	MealID  int64
	Reviews []models.Review
	Form    reviewForm
}

// mealPage is the data for the meal page.
type mealPage struct {
	Meal    models.MealDetail
	Upvote  upvoteButton
	Reviews reviewsSection
}

// meal shows one meal with its upvote button and reviews: GET /meals/{id}.
func (h *Handler) meal(w http.ResponseWriter, r *http.Request) {
	meal, ok := h.lookupMeal(w, r)
	if !ok {
		return
	}
	page, err := h.mealPageData(r.Context(), meal, visitorID(r), reviewForm{MealID: meal.ID})
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "meal", page)
}

// upvote toggles the visitor's upvote: POST /meals/{id}/upvote.
// HTMX gets back just the new button; a plain form post is redirected to
// the meal page (the Post/Redirect/Get pattern, so refresh won't re-submit).
func (h *Handler) upvote(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		h.notFound(w, r)
		return
	}
	count, upvoted, err := h.store.ToggleUpvote(r.Context(), id, visitorID(r))
	if errors.Is(err, store.ErrNotFound) {
		h.notFound(w, r)
		return
	}
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	if !isFragmentRequest(r) {
		http.Redirect(w, r, fmt.Sprintf("/meals/%d", id), http.StatusSeeOther)
		return
	}
	h.renderPartial(w, r, http.StatusOK, "meal", "upvote_button", upvoteButton{MealID: id, Count: count, Upvoted: upvoted})
}

// createReview saves a review: POST /meals/{id}/reviews.
func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	meal, ok := h.lookupMeal(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxReviewFormBytes)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "That review was too large to submit.", http.StatusRequestEntityTooLarge)
		return
	}
	// Keep exactly what they typed, so an error doesn't wipe their text.
	form := reviewForm{
		MealID: meal.ID,
		Author: r.PostForm.Get("author_name"),
		Body:   r.PostForm.Get("body"),
	}

	// Validate before rate limiting, so typos don't use up their 5 posts.
	author, body, err := store.NormalizeReview(form.Author, form.Body)
	var inputErr *store.ReviewInputError
	if errors.As(err, &inputErr) {
		form.AuthorErr, form.BodyErr = inputErr.Author, inputErr.Body
		// 422 = "we understood the request, but the data isn't valid".
		h.respondReviewForm(w, r, meal, form, http.StatusUnprocessableEntity)
		return
	}

	if !h.reviewLimiter.Allow(visitorID(r)) {
		form.FormErr = fmt.Sprintf("You've posted %d reviews in the last hour. Please wait a bit before posting another.", reviewLimit)
		h.respondReviewForm(w, r, meal, form, http.StatusTooManyRequests)
		return
	}

	if _, err := h.store.CreateReview(r.Context(), meal.ID, author, body); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.notFound(w, r)
			return
		}
		h.serverError(w, r, err)
		return
	}

	if !isFragmentRequest(r) {
		http.Redirect(w, r, fmt.Sprintf("/meals/%d#reviews", meal.ID), http.StatusSeeOther)
		return
	}
	// A fresh, empty form (so it "clears") above the updated list.
	h.respondReviewForm(w, r, meal, reviewForm{MealID: meal.ID, Posted: true}, http.StatusOK)
}

// respondReviewForm sends the reviews section with the given form state:
// just the section for HTMX, or the whole meal page for plain form posts.
func (h *Handler) respondReviewForm(w http.ResponseWriter, r *http.Request, meal models.MealDetail, form reviewForm, status int) {
	if isFragmentRequest(r) {
		section, err := h.reviewsSectionData(r.Context(), meal.ID, form)
		if err != nil {
			h.serverError(w, r, err)
			return
		}
		h.renderPartial(w, r, status, "meal", "reviews_section", section)
		return
	}
	page, err := h.mealPageData(r.Context(), meal, visitorID(r), form)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	h.render(w, r, status, "meal", page)
}

// lookupMeal loads the meal named by the {id} in the URL. If it can't, it
// writes the 404/500 response itself and returns ok=false.
func (h *Handler) lookupMeal(w http.ResponseWriter, r *http.Request) (models.MealDetail, bool) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		h.notFound(w, r)
		return models.MealDetail{}, false
	}
	meal, err := h.store.GetMeal(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		h.notFound(w, r)
		return models.MealDetail{}, false
	}
	if err != nil {
		h.serverError(w, r, err)
		return models.MealDetail{}, false
	}
	return meal, true
}

func (h *Handler) mealPageData(ctx context.Context, meal models.MealDetail, visitor string, form reviewForm) (mealPage, error) {
	stats, err := h.store.MealStats(ctx, []int64{meal.ID}, visitor)
	if err != nil {
		return mealPage{}, err
	}
	section, err := h.reviewsSectionData(ctx, meal.ID, form)
	if err != nil {
		return mealPage{}, err
	}
	st := stats[meal.ID]
	return mealPage{
		Meal:    meal,
		Upvote:  upvoteButton{MealID: meal.ID, Count: st.Upvotes, Upvoted: st.Upvoted},
		Reviews: section,
	}, nil
}

func (h *Handler) reviewsSectionData(ctx context.Context, mealID int64, form reviewForm) (reviewsSection, error) {
	reviews, err := h.store.ListReviewsByMeal(ctx, mealID)
	if err != nil {
		return reviewsSection{}, err
	}
	return reviewsSection{MealID: mealID, Reviews: reviews, Form: form}, nil
}

// mealCards pairs each meal with its upvote/review stats for this visitor.
func (h *Handler) mealCards(ctx context.Context, meals []models.Meal, visitor string) ([]mealCard, error) {
	if len(meals) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(meals))
	for i, m := range meals {
		ids[i] = m.ID
	}
	stats, err := h.store.MealStats(ctx, ids, visitor)
	if err != nil {
		return nil, err
	}
	cards := make([]mealCard, len(meals))
	for i, m := range meals {
		st := stats[m.ID]
		cards[i] = mealCard{
			Meal:        m,
			Upvote:      upvoteButton{MealID: m.ID, Count: st.Upvotes, Upvoted: st.Upvoted},
			ReviewCount: st.Reviews,
		}
	}
	return cards, nil
}
