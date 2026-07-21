package assets

import (
	"embed"
)

//go:embed "emails" "templates" all:static "models"
var EmbeddedFiles embed.FS
