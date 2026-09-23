# Slice 2 — Tap a Restaurant, See Its Meals

## Goal
A user taps a restaurant and sees its gluten-free meals, each with a name, photo, and short description.

## User story
> As a gluten-free diner, I want to see what GF meals a restaurant offers, so I can decide what to order before I go.

## What to build
1. **Store functions:**
   - `GetRestaurant(ctx, id)` — include the city name.
   - `ListMealsByRestaurant(ctx, restaurantID)` — sorted by name.
2. **Restaurant page (`GET /restaurants/{id}`):**
   - Top section: restaurant name, cuisine, address, description.
   - "← Back to {City} restaurants" link that returns to `/restaurants?city={slug}`.
   - Meal list below it.
3. **Meal card** shows: photo, meal name, short description.
   - Images use `loading="lazy"`, a fixed aspect ratio (so the page doesn't jump while loading), and `alt` text set to the meal name.
   - If `photo_url` is empty, show a default placeholder image.
4. **Edge cases:**
   - Restaurant with no meals → "No GF meals listed yet."
   - Non-number or unknown id → 404 page.

## Not in this slice
Reviews, upvotes, meal detail page, allergen charts, substitution info, GF safety Q&A.

## Done when
- [ ] Tapping a restaurant card from Slice 1 opens its page.
- [ ] Each meal shows photo, name, and description.
- [ ] Missing photos show the placeholder.
- [ ] Back link returns to the correct city list.
- [ ] Empty and unknown-restaurant cases behave as described.
- [ ] Works and looks good on a phone-sized screen.
- [ ] Tests: store test for `ListMealsByRestaurant`, handler tests for found, not-found, and bad id.
