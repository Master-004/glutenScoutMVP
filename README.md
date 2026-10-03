# GlutenScout

Helps gluten-free diners find restaurants with safe food options.
Go + `net/http` + `html/template` + HTMX + PostgreSQL.

> Sample data is **fictional**. Always confirm with a restaurant before ordering.

## Run it

You need [Docker Desktop](https://www.docker.com/products/docker-desktop/) running.

```sh
docker compose up --build
```

This starts everything in order: **db** (Postgres 16) → **migrate** (creates tables) →
**seed** (loads sample data) → **app** (web server). Then open <http://localhost:8080>.

- Stop: `Ctrl+C`, or `docker compose down`
- Start over with an empty database: `docker compose down -v` (deletes the data volume)
- After changing Go code or templates, run `docker compose up --build` again

Seeding is safe to repeat. It never duplicates rows or wipes data.

## Run the tests

You need Go installed. The version comes from `go.mod`, and Go downloads it automatically if needed.

```sh
docker compose up -d db   # store tests need the database
go test ./...
```

Handler tests use a fake store and run without Docker. Store tests use a separate
`glutenscout_test` database and **skip** unless `TEST_DATABASE_URL` is set. On Windows
PowerShell:

```powershell
$env:TEST_DATABASE_URL = "postgres://glutenscout:glutenscout@localhost:5432/glutenscout_test?sslmode=disable"
go test ./...
```

## Makefile (optional)

If you have `make`, these targets wrap the same commands: `make db-up`, `make migrate`,
`make seed`, `make run`, `make test`, `make down`.

## Project layout

```
cmd/server/       web server entry point
cmd/migrate/      applies migrations (run by Docker Compose)
internal/store/   all database queries
internal/handlers/ HTTP handlers + template rendering
internal/models/  plain structs (City, Restaurant, Meal)
migrations/       goose SQL migrations (embedded in the binary)
seed/seed.sql     fictional sample data
web/templates/    layout, pages/, partials/
web/static/       css, htmx, images (embedded in the binary)
docs/             project guide and slice specs
```

## Pages and routes

| Route | What it does |
|---|---|
| `GET /` | Home: pick a city |
| `GET /restaurants?city={slug}&type={salad\|grain\|entree\|dessert}` | A city's restaurants (full page, or just the list for HTMX); `type` keeps restaurants with a GF option of that type |
| `GET /restaurants/{id}?type=…` | A restaurant: food rating, photos, safety profile, call-ahead sheet, and its meals grouped by gluten status (`type` filters them) |
| `GET /meals/{id}` | A meal: photo, upvote button, reviews, review form |
| `POST /meals/{id}/upvote` | Toggle this visitor's upvote (returns the new button for HTMX) |
| `POST /meals/{id}/reviews` | Post a review (returns the updated reviews section for HTMX) |
| `POST /theme` | Save Auto / Light / Dark in a cookie, then go back to the previous page |
| `GET /healthz` | Returns `ok` when the database is reachable |

## How a few things work

- **No accounts (yet).** Each browser gets a random anonymous visitor ID in an
  `HttpOnly`, `SameSite=Lax` cookie. It allows one upvote per meal per visitor and
  limits reviews to 5 per visitor per hour (in memory, so the limit resets on restart).
- **HTMX, no custom JavaScript.** Interactivity is `hx-*` attributes in the templates.
  Every form also works as a plain HTML form.
- **Safety.** Review text is shown as plain text (`html/template` escapes it).
  Cross-site POSTs are blocked with Go's `http.CrossOriginProtection`, and all SQL
  uses query parameters.
