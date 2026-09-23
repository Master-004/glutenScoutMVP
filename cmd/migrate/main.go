// Command migrate creates/updates the database tables by applying the SQL
// files in migrations/. Docker Compose runs it before the web server starts.
package main

import (
	"context"
	"log"
	"os"

	"github.com/Master-004/glutenscout/internal/store"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if err := store.Migrate(context.Background(), dbURL); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations up to date")
}
