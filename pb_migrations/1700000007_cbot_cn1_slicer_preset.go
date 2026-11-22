package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		rec, err := app.FindFirstRecordByFilter("slicer_settings", "active = true")
		if err != nil || rec == nil {
			return nil
		}
		if rec.GetString("name") == "Jiddat CoreXY / Mark-1 0.20mm Production Profile" {
			rec.Set("name", "Cbot CN-1 0.4mm Nozzle (V11 Production Profile)")
			return app.Save(rec)
		}
		return nil
	}, func(app core.App) error {
		return nil
	})
}
