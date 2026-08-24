package migrations

import "embed"

// FS innehåller alla SQL-migrationer inbäddade i binären vid kompilering.
// Används av main.go via golang-migrate:s iofs-driver.
//
//go:embed *.sql
var FS embed.FS
