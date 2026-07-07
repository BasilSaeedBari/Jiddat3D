package pb_migrations

import (
	"encoding/json"
	
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

var collectionsJson = `[
	{
		"id": "collection_content",
		"name": "content",
		"type": "base",
		"system": false,
		"schema": [
			{"id": "title", "name": "title", "type": "text", "required": true, "options": {}},
			{"id": "slug", "name": "slug", "type": "text", "required": true, "options": {"pattern": "^[a-z0-9-]+$"}},
			{"id": "type", "name": "type", "type": "select", "required": true, "options": {"values": ["learn", "blog", "project"], "maxSelect": 1}},
			{"id": "track", "name": "track", "type": "select", "options": {"values": ["getting-started", "designing-making", "going-deeper"], "maxSelect": 1}},
			{"id": "excerpt", "name": "excerpt", "type": "text", "required": true, "options": {}},
			{"id": "body", "name": "body", "type": "editor", "required": true, "options": {}},
			{"id": "hero_image", "name": "hero_image", "type": "file", "required": true, "options": {"maxSelect": 1, "mimeTypes": ["image/jpeg", "image/png", "image/webp"]}},
			{"id": "gallery_images", "name": "gallery_images", "type": "file", "options": {"maxSelect": 99, "mimeTypes": ["image/jpeg", "image/png", "image/webp"]}},
			{"id": "attachment_files", "name": "attachment_files", "type": "file", "options": {"maxSelect": 99}},
			{"id": "author_name", "name": "author_name", "type": "text", "options": {}},
			{"id": "tags", "name": "tags", "type": "select", "options": {"values": ["hardware", "software", "community", "update"], "maxSelect": 10}},
			{"id": "has_math", "name": "has_math", "type": "bool", "options": {}},
			{"id": "featured", "name": "featured", "type": "bool", "options": {}},
			{"id": "published", "name": "published", "type": "bool", "options": {}},
			{"id": "published_at", "name": "published_at", "type": "date", "options": {}},
			{"id": "seo_description", "name": "seo_description", "type": "text", "options": {}}
		],
		"listRule": "",
		"viewRule": "",
		"createRule": null,
		"updateRule": null,
		"deleteRule": null,
		"indexes": ["CREATE UNIQUE INDEX idx_content_slug ON content (slug)"]
	},
	{
		"id": "collection_products",
		"name": "products",
		"type": "base",
		"system": false,
		"schema": [
			{"id": "name", "name": "name", "type": "text", "required": true, "options": {}},
			{"id": "slug", "name": "slug", "type": "text", "required": true, "options": {"pattern": "^[a-z0-9-]+$"}},
			{"id": "category", "name": "category", "type": "select", "required": true, "options": {"values": ["3d-printer", "cnc", "laser", "accessory"], "maxSelect": 1}},
			{"id": "tagline", "name": "tagline", "type": "text", "required": true, "options": {}},
			{"id": "description", "name": "description", "type": "editor", "options": {}},
			{"id": "specs", "name": "specs", "type": "json", "options": {}},
			{"id": "price", "name": "price", "type": "number", "options": {}},
			{"id": "currency", "name": "currency", "type": "text", "options": {}},
			{"id": "availability", "name": "availability", "type": "select", "required": true, "options": {"values": ["in-stock", "pre-order", "coming-soon"], "maxSelect": 1}},
			{"id": "hero_image", "name": "hero_image", "type": "file", "required": true, "options": {"maxSelect": 1, "mimeTypes": ["image/jpeg", "image/png", "image/webp"]}},
			{"id": "gallery_images", "name": "gallery_images", "type": "file", "options": {"maxSelect": 99, "mimeTypes": ["image/jpeg", "image/png", "image/webp"]}},
			{"id": "featured", "name": "featured", "type": "bool", "options": {}},
			{"id": "active", "name": "active", "type": "bool", "options": {}},
			{"id": "whatsapp_message_template", "name": "whatsapp_message_template", "type": "text", "options": {}}
		],
		"listRule": "",
		"viewRule": "",
		"createRule": null,
		"updateRule": null,
		"deleteRule": null,
		"indexes": ["CREATE UNIQUE INDEX idx_products_slug ON products (slug)"]
	},
	{
		"id": "collection_subscribers",
		"name": "subscribers",
		"type": "base",
		"system": false,
		"schema": [
			{"id": "email", "name": "email", "type": "email", "required": true, "options": {}},
			{"id": "active", "name": "active", "type": "bool", "options": {}},
			{"id": "source", "name": "source", "type": "text", "options": {}},
			{"id": "subscribed_at", "name": "subscribed_at", "type": "date", "options": {}}
		],
		"listRule": null,
		"viewRule": null,
		"createRule": "",
		"updateRule": null,
		"deleteRule": null,
		"indexes": ["CREATE UNIQUE INDEX idx_subscribers_email ON subscribers (email)"]
	},
	{
		"id": "collection_contact",
		"name": "contact_submissions",
		"type": "base",
		"system": false,
		"schema": [
			{"id": "name", "name": "name", "type": "text", "required": true, "options": {}},
			{"id": "contact_method", "name": "contact_method", "type": "text", "required": true, "options": {}},
			{"id": "subject", "name": "subject", "type": "select", "options": {"values": ["general", "product-inquiry", "support", "partnership"], "maxSelect": 1}},
			{"id": "message", "name": "message", "type": "text", "required": true, "options": {}},
			{"id": "related_product", "name": "related_product", "type": "relation", "options": {"collectionId": "collection_products", "cascadeDelete": false, "maxSelect": 1}},
			{"id": "status", "name": "status", "type": "select", "options": {"values": ["new", "read", "responded"], "maxSelect": 1}},
			{"id": "submitted_at", "name": "submitted_at", "type": "date", "options": {}}
		],
		"listRule": null,
		"viewRule": null,
		"createRule": "",
		"updateRule": null,
		"deleteRule": null,
		"indexes": []
	},
	{
		"id": "collection_settings",
		"name": "site_settings",
		"type": "base",
		"system": false,
		"schema": [
			{"id": "hero_headline", "name": "hero_headline", "type": "text", "options": {}},
			{"id": "hero_subheadline", "name": "hero_subheadline", "type": "text", "options": {}},
			{"id": "whatsapp_business_number", "name": "whatsapp_business_number", "type": "text", "options": {}},
			{"id": "instagram_url", "name": "instagram_url", "type": "url", "options": {}},
			{"id": "facebook_url", "name": "facebook_url", "type": "url", "options": {}},
			{"id": "youtube_url", "name": "youtube_url", "type": "url", "options": {}},
			{"id": "discord_url", "name": "discord_url", "type": "url", "options": {}},
			{"id": "impact_stat_1_label", "name": "impact_stat_1_label", "type": "text", "options": {}},
			{"id": "impact_stat_1_value", "name": "impact_stat_1_value", "type": "text", "options": {}},
			{"id": "impact_stat_2_label", "name": "impact_stat_2_label", "type": "text", "options": {}},
			{"id": "impact_stat_2_value", "name": "impact_stat_2_value", "type": "text", "options": {}},
			{"id": "impact_stat_3_label", "name": "impact_stat_3_label", "type": "text", "options": {}},
			{"id": "impact_stat_3_value", "name": "impact_stat_3_value", "type": "text", "options": {}},
			{"id": "impact_stat_4_label", "name": "impact_stat_4_label", "type": "text", "options": {}},
			{"id": "impact_stat_4_value", "name": "impact_stat_4_value", "type": "text", "options": {}},
			{"id": "manifesto_intro", "name": "manifesto_intro", "type": "editor", "options": {}}
		],
		"listRule": "",
		"viewRule": "",
		"createRule": null,
		"updateRule": null,
		"deleteRule": null,
		"indexes": []
	}
]`

func init() {
	m.Register(func(app core.App) error {
		var collections []map[string]any
		if err := json.Unmarshal([]byte(collectionsJson), &collections); err != nil {
			return err
		}

		return app.ImportCollections(collections, false)
	}, func(app core.App) error {
		return nil
	})
}
