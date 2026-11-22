package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Backfills products select options and missing seeded machines for databases
// that already ran 1700000004_business_seed before the schema fix.
func init() {
	m.Register(func(app core.App) error {
		return ensureProductsSchemaAndSeed(app)
	}, func(app core.App) error {
		return nil
	})
}
