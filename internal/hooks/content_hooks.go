package hooks

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	"jiddat3d/internal/compiler"
)

// RegisterContentHooks registers the hooks for compiling markdown and caching HTML
func RegisterContentHooks(app core.App) {
	app.OnRecordAfterCreateSuccess("content").BindFunc(func(e *core.RecordEvent) error {
		return processContentRecord(e.Record)
	})

	app.OnRecordAfterUpdateSuccess("content").BindFunc(func(e *core.RecordEvent) error {
		return processContentRecord(e.Record)
	})

	app.OnRecordAfterDeleteSuccess("content").BindFunc(func(e *core.RecordEvent) error {
		slug := e.Record.GetString("slug")
		return compiler.DeleteStaticPage("content", slug)
	})
	
	app.OnRecordAfterCreateSuccess("products").BindFunc(func(e *core.RecordEvent) error {
		return processProductRecord(e.Record)
	})

	app.OnRecordAfterUpdateSuccess("products").BindFunc(func(e *core.RecordEvent) error {
		return processProductRecord(e.Record)
	})

	app.OnRecordAfterDeleteSuccess("products").BindFunc(func(e *core.RecordEvent) error {
		slug := e.Record.GetString("slug")
		return compiler.DeleteStaticPage("products", slug)
	})
}

func processContentRecord(record *core.Record) error {
	published := record.GetBool("published")
	slug := record.GetString("slug")
	
	if !published {
		// Try deleting it if it was unpublished
		_ = compiler.DeleteStaticPage("content", slug)
		return nil
	}

	rawBody := record.GetString("body")
	htmlBody, err := compiler.RenderMarkdown(rawBody)
	if err != nil {
		log.Println("Error compiling markdown:", err)
		return err
	}

	data := compiler.TemplateData{
		Title:       record.GetString("title"),
		Description: record.GetString("seo_description"),
		Body:        compiler.SafeHTML(htmlBody), // We'll add SafeHTML helper
		Record:      recordToMap(record),
	}

	pageType := record.GetString("type")
	templateName := pageType + "_detail.html"

	fullHTML, err := compiler.RenderTemplate(templateName, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		return err
	}

	return compiler.WriteStaticPage("content", slug, fullHTML)
}

func processProductRecord(record *core.Record) error {
	active := record.GetBool("active")
	slug := record.GetString("slug")
	
	if !active {
		_ = compiler.DeleteStaticPage("products", slug)
		return nil
	}

	rawBody := record.GetString("description")
	htmlBody, _ := compiler.RenderMarkdown(rawBody)

	data := compiler.TemplateData{
		Title:       record.GetString("name"),
		Description: record.GetString("tagline"),
		Body:        compiler.SafeHTML(htmlBody),
		Record:      recordToMap(record),
	}

	fullHTML, err := compiler.RenderTemplate("product_detail.html", data)
	if err != nil {
		log.Println("Error rendering product template:", err)
		return err
	}

	return compiler.WriteStaticPage("products", slug, fullHTML)
}

func recordToMap(r *core.Record) map[string]any {
	return r.PublicExport()
}
