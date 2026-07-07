package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 1. Content
		content := core.NewBaseCollection("content")
		content.Name = "content"
		content.Fields.Add(&core.TextField{Name: "title", Required: true})
		content.Fields.Add(&core.TextField{Name: "slug", Required: true})
		content.Fields.Add(&core.SelectField{Name: "type", Values: []string{"learn", "blog", "project"}, MaxSelect: 1, Required: true})
		content.Fields.Add(&core.SelectField{Name: "track", Values: []string{"getting-started", "designing-making", "going-deeper"}, MaxSelect: 1})
		content.Fields.Add(&core.TextField{Name: "excerpt", Required: true})
		content.Fields.Add(&core.EditorField{Name: "body", Required: true})
		content.Fields.Add(&core.FileField{Name: "hero_image", MaxSelect: 1, MimeTypes: []string{"image/jpeg", "image/png", "image/webp"}})
		content.Fields.Add(&core.FileField{Name: "gallery_images", MaxSelect: 99, MimeTypes: []string{"image/jpeg", "image/png", "image/webp"}})
		content.Fields.Add(&core.FileField{Name: "attachment_files", MaxSelect: 99})
		content.Fields.Add(&core.TextField{Name: "author_name"})
		content.Fields.Add(&core.SelectField{Name: "tags", Values: []string{"hardware", "software", "community", "update"}, MaxSelect: 4})
		content.Fields.Add(&core.BoolField{Name: "has_math"})
		content.Fields.Add(&core.BoolField{Name: "featured"})
		content.Fields.Add(&core.BoolField{Name: "published"})
		content.Fields.Add(&core.DateField{Name: "published_at"})
		content.Fields.Add(&core.TextField{Name: "seo_description"})
		content.AddIndex("idx_content_slug", false, "slug", "")
		if err := app.Save(content); err != nil { return err }

		// 2. Products
		products := core.NewBaseCollection("products")
		products.Name = "products"
		products.Fields.Add(&core.TextField{Name: "name", Required: true})
		products.Fields.Add(&core.TextField{Name: "slug", Required: true})
		products.Fields.Add(&core.SelectField{Name: "category", Values: []string{"3d-printer", "cnc", "laser", "accessory"}, MaxSelect: 1, Required: true})
		products.Fields.Add(&core.TextField{Name: "tagline", Required: true})
		products.Fields.Add(&core.EditorField{Name: "description"})
		products.Fields.Add(&core.JSONField{Name: "specs"})
		products.Fields.Add(&core.NumberField{Name: "price"})
		products.Fields.Add(&core.TextField{Name: "currency"})
		products.Fields.Add(&core.SelectField{Name: "availability", Values: []string{"in-stock", "pre-order", "coming-soon"}, MaxSelect: 1, Required: true})
		products.Fields.Add(&core.FileField{Name: "hero_image", MaxSelect: 1, MimeTypes: []string{"image/jpeg", "image/png", "image/webp"}})
		products.Fields.Add(&core.FileField{Name: "gallery_images", MaxSelect: 99, MimeTypes: []string{"image/jpeg", "image/png", "image/webp"}})
		products.Fields.Add(&core.BoolField{Name: "featured"})
		products.Fields.Add(&core.BoolField{Name: "active"})
		products.Fields.Add(&core.TextField{Name: "whatsapp_message_template"})
		products.AddIndex("idx_products_slug", false, "slug", "")
		if err := app.Save(products); err != nil { return err }

		// 3. Subscribers
		subscribers := core.NewBaseCollection("subscribers")
		subscribers.Name = "subscribers"
		subscribers.Fields.Add(&core.EmailField{Name: "email", Required: true})
		subscribers.Fields.Add(&core.BoolField{Name: "active"})
		subscribers.Fields.Add(&core.TextField{Name: "source"})
		subscribers.Fields.Add(&core.DateField{Name: "subscribed_at"})
		subscribers.AddIndex("idx_subscribers_email", true, "email", "")
		if err := app.Save(subscribers); err != nil { return err }

		// 4. Contact Submissions
		contact := core.NewBaseCollection("contact_submissions")
		contact.Name = "contact_submissions"
		contact.Fields.Add(&core.TextField{Name: "name", Required: true})
		contact.Fields.Add(&core.TextField{Name: "contact_method", Required: true})
		contact.Fields.Add(&core.SelectField{Name: "subject", Values: []string{"general", "product-inquiry", "support", "partnership"}, MaxSelect: 1})
		contact.Fields.Add(&core.TextField{Name: "message", Required: true})
		contact.Fields.Add(&core.RelationField{Name: "related_product", CollectionId: products.Id, MaxSelect: 1})
		contact.Fields.Add(&core.SelectField{Name: "status", Values: []string{"new", "read", "responded"}, MaxSelect: 1})
		contact.Fields.Add(&core.DateField{Name: "submitted_at"})
		if err := app.Save(contact); err != nil { return err }

		// 5. Site Settings
		settings := core.NewBaseCollection("site_settings")
		settings.Name = "site_settings"
		settings.Fields.Add(&core.TextField{Name: "hero_headline"})
		settings.Fields.Add(&core.TextField{Name: "hero_subheadline"})
		settings.Fields.Add(&core.TextField{Name: "whatsapp_business_number"})
		settings.Fields.Add(&core.URLField{Name: "instagram_url"})
		settings.Fields.Add(&core.URLField{Name: "facebook_url"})
		settings.Fields.Add(&core.URLField{Name: "youtube_url"})
		settings.Fields.Add(&core.URLField{Name: "discord_url"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_1_label"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_1_value"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_2_label"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_2_value"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_3_label"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_3_value"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_4_label"})
		settings.Fields.Add(&core.TextField{Name: "impact_stat_4_value"})
		settings.Fields.Add(&core.EditorField{Name: "manifesto_intro"})
		if err := app.Save(settings); err != nil { return err }

		return nil
	}, func(app core.App) error {
		return nil
	})
}
