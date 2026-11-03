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
			settingsRecord.Set("hero_headline", "Engineering Innovation in Pakistan. Specialized 3D Printers.")
			settingsRecord.Set("hero_subheadline", "Custom-designed FDM/FFM 3D printers, 100% open source architecture, 6-Month Working & Education Warranty. Precision-built in Karachi on a prepaid basis.")
			_ = app.Save(settingsRecord)
		}

		// 2. Seed products with our own specialized Jiddat machines (No ready-made international imports)
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
				name:     "Jiddat Apex CoreXY",
				slug:     "jiddat-apex-corexy",
				category: "3d-printer",
				tagline:  "Industrial enclosed high-speed FDM/FFM workhorse running custom Klipper firmware. Built for engineering plastics (ABS, Nylon, PC).",
				price:    245000,
				avail:    "built-to-order",
				featured: true,
			},
			{
				name:     "Jiddat Mark-1 Desktop",
				slug:     "jiddat-mark-1",
				category: "3d-printer",
				tagline:  "Precision direct-drive FDM/FFM printer with CNC-machined rigid chassis and all-metal hotend for makers and engineers.",
				price:    145000,
				avail:    "in-stock",
				featured: true,
			},
			{
				name:     "Jiddat Edu-Maker V2",
				slug:     "jiddat-edu-maker",
				category: "3d-printer",
				tagline:  "Classroom-safe, rugged educational FDM/FFM printer designed for Pakistani schools, colleges, and university STEAM labs.",
				price:    98000,
				avail:    "in-stock",
				featured: true,
			},
			{
				name:     "Jiddat L-500 Large-Format",
				slug:     "jiddat-l500-workstation",
				category: "3d-printer",
				tagline:  "Massive 500x500x500mm build chamber for industrial tooling, full-scale automotive parts, and final year project prototypes.",
				price:    385000,
				avail:    "built-to-order",
				featured: true,
			},
			{
				name:     "Jiddat CoreXY Maker Edition",
				slug:     "jiddat-corexy",
				category: "3d-printer",
				tagline:  "100% open-source RepRap DIY maker assembly kit with local precision hardware, open firmware, and full CAD schematics.",
				price:    78000,
				avail:    "in-stock",
				featured: false,
			},
			{
				name:     "Jiddat Precision Resin (MSLA)",
				slug:     "jiddat-precision-resin",
				category: "resin-printer",
				tagline:  "Ultra-high resolution stereolithography system for jewelry and dental casting — Currently in Engineering R&D in Karachi.",
				price:    135000,
				avail:    "in-development",
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
