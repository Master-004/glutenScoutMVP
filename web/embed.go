// Package web bundles templates and static files into the Go binary with
// go:embed, so the server runs without needing these files on disk.
package web

import "embed"

//go:embed templates static
var FS embed.FS
