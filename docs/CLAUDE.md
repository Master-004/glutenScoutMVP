# GlutenScout — Project Guide for Claude Code

## What we're building
GlutenScout helps gluten-free (GF) people quickly find restaurants with safe food options.

**Main user flow:** Pick a location → see GF-friendly restaurants → tap a restaurant to see its meals → review/upvote meals they trust.

This is an MVP. Build only what the current slice asks for. Nice-to-have features are listed at the bottom so you know to leave them out.

## How we work
- The MVP is split into 4 slices. Each slice is a working, testable feature from database to screen.
- Build one slice at a time, in order: `slice-0-foundation.md` → `slice-1-location-and-restaurants.md` → `slice-2-restaurant-meals.md` → `slice-3-reviews-and-upvotes.md`.
- Before writing code for a slice, read its file and summarize your plan. Ask if anything is unclear.
- A slice is done only when every item in its "Done when" checklist passes.
- Keep things simple and readable. The developer is comfortable with Go but newer to web apps, so add short comments explaining *why*, not just *what*.

## Tech stack (decided — don't swap these out)
- **Language:** Go (1.22 or newer) for both backend and frontend.
- **Web server & routing:** Go standard library `net/http` (use the Go 1.22 `ServeMux` method/path patterns like `GET /restaurants/{id}`). No web framework.
- **Pages:** Server-rendered HTML with Go's `html/template`.
- **Interactivity:** HTMX (loaded from a local file in `web/static/js/`) for small partial-page updates. No React, no custom JavaScript unless a slice says so.
- **Styling:** One plain CSS file (`web/static/css/app.css`). Mobile-first — most users will be on phones.
- **Database:** PostgreSQL, run locally with Docker Compose.
- **DB driver:** `github.com/jackc/pgx/v5` (use `pgxpool`).
- **Migrations:** `github.com/pressly/goose/v3`, SQL files in `migrations/`.
- **Config:** Environment variables (e.g. `DATABASE_URL`, `PORT`) with a `.env.example` file.
- **Tests:** Go's built-in `testing` package and `net/http/httptest`.

## Project layout
```
glutenscout/
├── cmd/server/main.go        # starts the web server
├── internal/
│   ├── store/                # all database queries live here
│   ├── handlers/             # HTTP handlers (one file per feature area)
│   └── models/               # plain Go structs (City, Restaurant, Meal, ...)
├── migrations/               # goose SQL migrations
├── seed/seed.sql             # sample data
├── web/
│   ├── templates/            # layout.html + page and partial templates
│   └── static/               # css, js (htmx), images
├── docker-compose.yml
├── Makefile                  # make db-up, make migrate, make seed, make run, make test
└── .env.example
```

## Rules of the road
- Handlers stay thin: read the request, call `store`, render a template.
- All SQL lives in `internal/store`. Always use query parameters (`$1`, `$2`), never string-building SQL.
- Templates must escape user content (the default with `html/template` — don't use `template.HTML` on user input).
- Show friendly error pages (404 and 500), never raw Go errors, to users.
- Sample data uses **fictional** restaurants. Never claim a real restaurant is gluten-free.
- Every page footer shows this disclaimer: *"GlutenScout info comes from the community and may be out of date. Always confirm with the restaurant before ordering."*

## Added after the 4 slices (follow docs/DESIGN.md)
- **Safety profile** on the restaurant page: "Last verified" date, precaution badges (confirmed / conflict / unverified), restaurant-answers date, conflict note in words, and a call-ahead script with the restaurant's phone. All of it comes from **seed data** (`safety_profiles`, `precautions`) for now.
- **Reviews** require a reviewer type (Celiac / Strict GF / Gluten-sensitive / Restaurant) and a 1–5 **food** rating (restaurants don't rate their own food). Average ratings show on restaurant rows, restaurant pages, meal cards and meal pages.
- **Theme switch** (Auto / Light / Dark) saved in a cookie.
- **Meal gluten status** (`meals.gf_status`: `gf` / `cross_contact` / `contains`, plus `gf_note`). The restaurant page groups meals: GF options, then cross-contact risk, then "Contains gluten" folded in a `<details>`. "N GF meals" counts only `gf`. Non-GF meals are in the seed on purpose, to test whether diners can find the GF options.
- **Sheets**: "Call ahead" and "Write a review" open bottom sheets built on the HTML `popover` attribute (no JavaScript). Review errors with HTMX swap only the form (`HX-Retarget`) so the sheet stays open.
- **Green button**: "Write a review" on the meal page. List and restaurant pages have no green page button (the cards are the main action).

- **Round 1 prototype tweaks** (docs/PROTOTYPE_TWEAKS.md): "Find restaurants" button and next-step hints; "No reviews yet" states; safety profile explainer + info sheet; restaurant photos (Food / Storefront / Prep area, seed illustrations, "Add a photo" = coming soon); **food-type filter chips** (`meals.food_type`; `?type=` on the restaurant list and restaurant page; the list keeps restaurants with a *GF option* of that type); "Good for mixed groups" (`restaurants.good_for_groups`).

## NOT in the MVP (don't build these yet)
User accounts/login, verified users, real photo uploads, sorting and other filters (distance, best reviews, most options, reliability), allergen charts, other allergens (dairy, nuts, etc.), user-submitted precaution reports (the safety profile is seed data only and is not computed from reviews yet), peer review actions (Helpful / Matches my visit / Report conflict), substitution info, community events, payments, a Swift/iOS app, and outside restaurant APIs (Google Places, Foursquare, etc.).
