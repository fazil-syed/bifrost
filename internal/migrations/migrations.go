package migrations

import (
	"embed"
)

//go:embed global/*.sql
var Global embed.FS
