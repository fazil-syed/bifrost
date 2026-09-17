package migrations

import (
	"embed"
)

//go:embed global/*.sql
var Global embed.FS

//go:embed tenant/*.sql
var Tenant embed.FS
