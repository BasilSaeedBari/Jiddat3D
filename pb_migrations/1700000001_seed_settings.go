package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)
		record.Set("hero_headline", "Build Something Amazing.")
		record.Set("hero_subheadline", "Affordable, repairable manufacturing tools — designed, assembled, and supported from right here in Pakistan.")
		record.Set("whatsapp_business_number", "+920000000000") // Placeholder
		
		record.Set("impact_stat_1_label", "Units in the field")
		record.Set("impact_stat_1_value", "0")
		record.Set("impact_stat_2_label", "Active makers")
		record.Set("impact_stat_2_value", "0")
		record.Set("impact_stat_3_label", "Projects shared")
		record.Set("impact_stat_3_value", "0")
		record.Set("impact_stat_4_label", "Guides published")
		record.Set("impact_stat_4_value", "0")

		record.Set("manifesto_intro", "I started building printers because I wanted affordable manufacturing in Pakistan. I wanted engineers, students, and independent makers to own real machines — not just admire them in imported catalogs. Jiddat means innovation — a deliberate departure from the old way of doing things. That's what we're building toward.")

		return app.Save(record)
	}, func(app core.App) error {
		return nil
	})
}
