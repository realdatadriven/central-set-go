package assets

import (
	"embed"
)

//go:embed "emails" "migrations" "setup" "templates" "static"
var EmbeddedFiles embed.FS
