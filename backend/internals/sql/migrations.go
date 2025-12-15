package sql

import "embed"

//go:embed schemas/*.sql
var EmbedMigrations embed.FS

//go:embed seed/*.sql
var EmbedSeed embed.FS
