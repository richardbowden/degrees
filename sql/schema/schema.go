package schema

import "embed"

//go:embed *.sql
var SQLMigrationFiles embed.FS
