package sql

import "embed"

//go:embed schemas/*.sql
var EmbedMigrations embed.FS
