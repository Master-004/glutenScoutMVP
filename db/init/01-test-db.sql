-- Runs once when the Postgres volume is first created.
-- A separate database for `go test` so tests can wipe tables freely.
CREATE DATABASE glutenscout_test;
