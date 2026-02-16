package assets

import (
	"embed"
)

//go:embed "emails" "migrations" "setup" "templates"
var EmbeddedFiles embed.FS
