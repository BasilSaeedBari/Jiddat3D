package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 1. Update site_settings
		settingsRecord, err := app.FindFirstRecordByFilter("site_settings", "1=1")
		if err == nil && settingsRecord != nil {
			settingsRecord.Set("whatsapp_business_number", "923001318112")
			settingsRecord.Set("hero_headline", "3D Printers for Pakistan. Bought Here. Fixed Here. Taught Here.")
			settingsRecord.Set("hero_subheadline", "Genuine machines with local warranty, spare parts, and hands-on training. Delivered across Pakistan with Cash on Delivery.")
			_ = app.Save(settingsRecord)
		}

		// 2. Seed products
		productsCol, err := app.FindCollectionByNameOrId("products")
		if err != nil {
			return nil
		}

		seedList := []struct {
			name        string
			slug        string
			category    string
			tagline     string
			price       float64
			avail       string
			featured    bool
		}{
			{
				name:     "Bambu Lab A1",
				slug:     "bambu-lab-a1",
				category: "3d-printer",
				tagline:  "Flagship automated 500mm/s high-speed printer with optional 4-color AMS Lite.",
				price:    165000,
				avail:    "in-stock",
				featured: true,
			},
			{
				name:     "Bambu Lab A1 Mini",
				slug:     "bambu-lab-a1-mini",
				category: "3d-printer",
				tagline:  "Compact, desktop automated 3D printer with active noise cancellation.",
				price:    115000,
				avail:    "in-stock",
				featured: true,
			},
			{
				name:     "Bambu Lab P1S Enclosed",
				slug:     "bambu-lab-p1s",
				category: "3d-printer",
				tagline:  "Fully enclosed CoreXY chamber designed for high-temperature ABS and engineering filaments.",
				price:    245000,
				avail:    "in-stock",
				featured: true,
			},
			{
				name:     "Creality Ender-3 V3 KE",
				slug:     "creality-ender-3-v3-ke",
				category: "3d-printer",
				tagline:  "High-speed 500mm/s CoreXZ architecture running Klipper OS under 1 Lakh PKR.",
				price:    98000,
				avail:    "in-stock",
				featured: false,
			},
			{
				name:     "Elegoo Mars 4 Ultra 9K",
				slug:     "elegoo-mars-4-ultra",
				category: "3d-printer",
				tagline:  "Ultra-precision 9K mono LCD resin printer for fine jewelry casting and dental molds.",
				price:    128000,
				avail:    "in-stock",
				featured: false,
			},
			{
				name:     "Jiddat CoreXY Maker Kit",
				slug:     "jiddat-corexy",
				category: "3d-printer",
				tagline:  "100% repairable open-source RepRap machine assembled with local Pakistani hardware.",
				price:    85000,
				avail:    "in-stock",
				featured: false,
			},
		}

		for _, item := range seedList {
			existing, _ := app.FindFirstRecordByFilter("products", "slug = {:slug}", map[string]any{"slug": item.slug})
			if existing != nil {
				continue
			}

			rec := core.NewRecord(productsCol)
			rec.Set("name", item.name)
			rec.Set("slug", item.slug)
			rec.Set("category", item.category)
			rec.Set("tagline", item.tagline)
			rec.Set("price", item.price)
			rec.Set("currency", "PKR")
			rec.Set("availability", item.avail)
			rec.Set("featured", item.featured)
			rec.Set("active", true)
			_ = app.Save(rec)
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}
