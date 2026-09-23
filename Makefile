# Thin shortcuts. You don't need make: `docker compose up --build` does it all.
.PHONY: db-up migrate seed run test down

db-up:
	docker compose up -d db

migrate:
	docker compose run --rm --build migrate

seed:
	docker compose run --rm seed

run:
	docker compose up --build

test:
	go test ./...

down:
	docker compose down
