package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 1. Slicer Settings Collection
		slicerSettings := core.NewBaseCollection("slicer_settings")
		slicerSettings.Name = "slicer_settings"
		slicerSettings.Fields.Add(&core.TextField{Name: "name", Required: true})
		slicerSettings.Fields.Add(&core.FileField{Name: "profile_file", MaxSelect: 1})
		slicerSettings.Fields.Add(&core.NumberField{Name: "time_multiplier", Required: true})
		slicerSettings.Fields.Add(&core.NumberField{Name: "weight_multiplier", Required: true})
		slicerSettings.Fields.Add(&core.NumberField{Name: "machine_hourly_rate", Required: true})
		slicerSettings.Fields.Add(&core.NumberField{Name: "labor_overhead"})
		slicerSettings.Fields.Add(&core.NumberField{Name: "electricity_hourly"})
		slicerSettings.Fields.Add(&core.NumberField{Name: "failure_rate_pct"})
		slicerSettings.Fields.Add(&core.NumberField{Name: "profit_multiplier_pct"})
		slicerSettings.Fields.Add(&core.NumberField{Name: "default_filament_price", Required: true})
		slicerSettings.Fields.Add(&core.BoolField{Name: "active"})

		if err := app.Save(slicerSettings); err != nil {
			return err
		}

		// Seed Default Slicer Setting Record
		defaultSetting := core.NewRecord(slicerSettings)
		defaultSetting.Set("name", "Jiddat CoreXY / Mark-1 0.20mm Production Profile")
		defaultSetting.Set("time_multiplier", 1.5)
		defaultSetting.Set("weight_multiplier", 1.5)
		defaultSetting.Set("machine_hourly_rate", 250)
		defaultSetting.Set("labor_overhead", 50)
		defaultSetting.Set("electricity_hourly", 30)
		defaultSetting.Set("failure_rate_pct", 5)
		defaultSetting.Set("profit_multiplier_pct", 25)
		defaultSetting.Set("default_filament_price", 8)
		defaultSetting.Set("active", true)
		if err := app.Save(defaultSetting); err != nil {
			return err
		}

		// 2. Filament Types Collection
		filamentTypes := core.NewBaseCollection("filament_types")
		filamentTypes.Name = "filament_types"
		filamentTypes.Fields.Add(&core.TextField{Name: "name", Required: true})
		filamentTypes.Fields.Add(&core.TextField{Name: "material", Required: true})
		filamentTypes.Fields.Add(&core.TextField{Name: "color_name", Required: true})
		filamentTypes.Fields.Add(&core.TextField{Name: "color_hex", Required: true})
		filamentTypes.Fields.Add(&core.NumberField{Name: "density_g_cm3", Required: true})
		filamentTypes.Fields.Add(&core.NumberField{Name: "price_per_gram", Required: true})
		filamentTypes.Fields.Add(&core.BoolField{Name: "is_default"})
		filamentTypes.Fields.Add(&core.BoolField{Name: "active"})

		if err := app.Save(filamentTypes); err != nil {
			return err
		}

		// Seed Default Filaments
		filaments := []struct {
			name      string
			material  string
			colorName string
			colorHex  string
			density   float64
			price     int64
			isDefault bool
		}{
			{"White Basic PLA", "PLA", "Signal White", "#FAF7F2", 1.24, 8, true},
			{"Black Tough PLA+", "PLA+", "Matte Black", "#1E1E1E", 1.25, 10, false},
			{"High-Temp PETG", "PETG", "Terracotta Amber", "#D9530F", 1.27, 12, false},
			{"TPU 95A Flexible", "TPU", "Craft Teal", "#17C3B2", 1.21, 18, false},
		}

		for _, f := range filaments {
			rec := core.NewRecord(filamentTypes)
			rec.Set("name", f.name)
			rec.Set("material", f.material)
			rec.Set("color_name", f.colorName)
			rec.Set("color_hex", f.colorHex)
			rec.Set("density_g_cm3", f.density)
			rec.Set("price_per_gram", f.price)
			rec.Set("is_default", f.isDefault)
			rec.Set("active", true)
			if err := app.Save(rec); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}
