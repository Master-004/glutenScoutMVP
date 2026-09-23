# Slice 3 — Review and Upvote Meals

## Goal
Users can leave a written review on a meal and upvote meals they trust. This is what makes GlutenScout community-powered.

## User story
> As a gluten-free diner, I want to upvote meals that were safe for me and read others' reviews, so the community knows which meals to trust.

## Important decision: no accounts yet
Login is not part of the MVP. Instead:
- Give each browser an **anonymous visitor ID**: a random value stored in an `HttpOnly`, `SameSite=Lax` cookie, created by middleware on first visit.
- Use it to allow **one upvote per meal per visitor** (tapping again removes the upvote).
- Reviews ask for an optional display name (default "Anonymous").

## What to build
1. **Migration `0002_reviews_upvotes.sql`:**
   - `reviews` — `id`, `meal_id` (FK), `author_name`, `body`, `created_at`
   - `meal_upvotes` — `meal_id` (FK), `visitor_id`, `created_at`, primary key (`meal_id`, `visitor_id`)
2. **Visitor ID middleware** as described above.
3. **Store functions:**
   - `GetMeal(ctx, id)` (with restaurant name)
   - `ListReviewsByMeal(ctx, mealID)` — newest first
   - `CreateReview(ctx, mealID, authorName, body)`
   - `ToggleUpvote(ctx, mealID, visitorID)` → returns new count and whether this visitor has upvoted
   - Upvote counts for a list of meals (for the Slice 2 page)
4. **Meal page (`GET /meals/{id}`):** photo, name, description, restaurant link, upvote button with count, review list, and review form.
5. **Upvote button** (`POST /meals/{id}/upvote`):
   - Appears on each meal card on the restaurant page and on the meal page.
   - Uses HTMX to swap just the button (new count + highlighted if upvoted).
6. **Review form** (`POST /meals/{id}/reviews`):
   - Fields: display name (optional, max 40 chars) and review (required, 1–1000 chars).
   - Validate on the server. On error, re-show the form with a clear message and keep what they typed.
   - On success, the new review appears at the top of the list (HTMX swap) and the form clears.
7. **Link meal cards** on the restaurant page to `/meals/{id}` and show "X reviews."
8. **Basic safety:**
   - Reviews are displayed as plain text (escaped by `html/template`).
   - Protect POST forms from cross-site requests (check the `Origin`/`Sec-Fetch-Site` header, or add a CSRF token).
   - Simple rate limit: max 5 reviews per visitor per hour (in-memory is fine for MVP).

## Not in this slice
User accounts, verified users, photo uploads, star ratings, editing/deleting reviews, moderation tools, "sort by best reviewed."

## Done when
- [ ] Tapping upvote increases the count; tapping again removes it; refreshing keeps the right state.
- [ ] A second browser (or incognito window) counts as a different visitor.
- [ ] Submitting a valid review shows it at the top without a full reload.
- [ ] Empty or too-long reviews show a friendly error and keep the typed text.
- [ ] Typing `<script>` in a review shows it as plain text, not code.
- [ ] Restaurant page shows upvote counts and review counts on each meal.
- [ ] Works and looks good on a phone-sized screen.
- [ ] Tests: `ToggleUpvote` (add, remove, one-per-visitor), `CreateReview` validation, handler tests for both POST routes.

## After this slice
The MVP is complete. 🎉 Good next slices to consider: user accounts, a restaurant API for real locations, sorting/filters, and GF safety Q&A.
