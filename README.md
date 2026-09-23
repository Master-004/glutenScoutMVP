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

## Useful URLs

- `/` — home
- `/healthz` — returns `ok` when the database is reachable
