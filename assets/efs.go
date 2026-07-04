package assets

import (
	"embed"
)

//go:embed "emails" "templates" "static" "models"
var EmbeddedFiles embed.FS
