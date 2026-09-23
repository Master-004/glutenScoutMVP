# Slice 1 — Pick a Location, See Restaurants

## Goal
A user picks their city and immediately sees a list of GF-friendly restaurants there.

## User story
> As a gluten-free diner, I want to choose my city and see restaurants with GF options, so I know where I can eat.

## What to build
1. **Store functions:**
   - `ListCities(ctx)` — all cities, sorted by name.
   - `GetCityBySlug(ctx, slug)`
   - `ListRestaurantsByCity(ctx, cityID)` — sorted by name; include a count of meals for each restaurant.
2. **Home page (`GET /`):** replace the placeholder with:
   - A short welcome line ("Find gluten-free-friendly restaurants near you").
   - A city dropdown listing all cities.
   - An empty results area below it.
3. **Restaurant list (`GET /restaurants?city={slug}`):**
   - With HTMX: changing the dropdown fetches just the list partial and swaps it into the results area (no full page reload).
   - Without HTMX (someone visits the URL directly): render the full page with the dropdown pre-selected and the list shown. Use the `HX-Request` header to decide which to render.
   - Update the browser URL (`hx-push-url`) so the link can be shared or bookmarked.
4. **Restaurant card** shows: name, cuisine, address, and "X GF meals." Each card links to `/restaurants/{id}` (the page is built in Slice 2 — a 404 for now is fine).
5. **Empty and error states:**
   - City with no restaurants → "No restaurants here yet."
   - Unknown city slug → 404 page.

## Not in this slice
Restaurant detail/meals page, sorting, filters, distance, maps, GPS location, outside APIs.

## Done when
- [ ] Choosing a city shows that city's restaurants without a full page reload.
- [ ] Visiting `/restaurants?city=<slug>` directly shows the full page with the list.
- [ ] Each card shows name, cuisine, address, and meal count.
- [ ] Empty city and unknown city behave as described.
- [ ] Works and looks good on a phone-sized screen.
- [ ] Tests: store tests for `ListRestaurantsByCity`, handler tests for the full page and the HTMX partial.
