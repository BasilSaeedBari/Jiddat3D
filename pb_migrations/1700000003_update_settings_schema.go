package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		settings, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}

		// Add whatsapp_community_url
		settings.Fields.Add(&core.URLField{Name: "whatsapp_community_url"})

		if err := app.Save(settings); err != nil {
			return err
		}

		// Now update the actual record
		record, err := app.FindFirstRecordByFilter("site_settings", "1=1")
		if err != nil {
			return err
		}

		record.Set("whatsapp_business_number", "923001318112")
		record.Set("whatsapp_community_url", "https://chat.whatsapp.com/Fgvte09elKICDNzlPVuJTe")
		record.Set("discord_url", "https://discord.gg/rvsaqF2Q7p")

		return app.Save(record)
	}, func(app core.App) error {
		return nil
	})
}
