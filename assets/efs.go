package assets

import (
	"embed"
)

//go:embed "emails" "migrations" "setup"
var EmbeddedFiles embed.FS
