// Package migrations embeds the goose SQL files so both the migrate command
// and the tests can apply them without needing the files on disk.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
