# Slice 0 — Foundation

## Goal
Set up the empty app so every later slice has a place to plug in. Think of it as setting up the kitchen before cooking: running server, database, base page layout, and sample data.

## What to build
1. **Project setup:** Go module (`github.com/<you>/glutenscout`), the folder layout from `CLAUDE.md`, `.gitignore`, `.env.example`.
2. **Database:** `docker-compose.yml` running PostgreSQL 16 with a named volume.
3. **Migration `0001_init.sql`** with these tables:
   - `cities` — `id`, `name`, `state`, `slug` (unique, e.g. `denver-co`)
   - `restaurants` — `id`, `city_id` (FK), `name`, `address`, `cuisine`, `description`, `created_at`
   - `meals` — `id`, `restaurant_id` (FK), `name`, `description`, `photo_url`, `created_at`
   - Add indexes on `restaurants.city_id` and `meals.restaurant_id`.
4. **Sample data (`seed/seed.sql`):** 3 cities, 3–5 fictional restaurants per city, 3–5 GF meals per restaurant. Meal photos point to local placeholder images in `web/static/images/meals/` (simple generated placeholders are fine).
5. **Models:** Go structs for `City`, `Restaurant`, `Meal` in `internal/models`.
6. **Store:** a `Store` type wrapping a `pgxpool.Pool`, created in `main.go` and passed to handlers.
7. **Web server:** `cmd/server/main.go` reads `PORT` and `DATABASE_URL`, serves `/static/`, and has a `GET /healthz` route that checks the DB and returns `ok`.
8. **Base layout template:** header with the GlutenScout name, a main content area, and the footer disclaimer from `CLAUDE.md`. Includes `app.css` and HTMX.
9. **Placeholder home page** at `GET /` that says "GlutenScout — coming soon" using the layout.
10. **Error pages:** simple 404 and 500 templates.
11. **Makefile targets:** `db-up`, `migrate`, `seed`, `run`, `test`.

## Not in this slice
Any real features (city picker, restaurant list, meals, reviews).

## Done when
- [ ] `make db-up && make migrate && make seed && make run` works from a fresh clone.
- [ ] Visiting `http://localhost:8080/` shows the placeholder page with header and disclaimer footer.
- [ ] `GET /healthz` returns `ok`.
- [ ] Visiting a made-up URL shows the friendly 404 page.
- [ ] The page looks fine on a phone-sized screen (check with browser dev tools).
- [ ] `make test` passes (at least one handler test for `/healthz` using `httptest`).
- [ ] A short `README.md` explains how to run the app.
