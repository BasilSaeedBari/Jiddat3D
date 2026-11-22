package routes

import (
	"fmt"
	"html"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"jiddat3d/internal/compiler"
	"jiddat3d/internal/mailer"
)

func RegisterPublicRoutes(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Wire up instant 3D model slicer and reactive quotation routes
		RegisterSlicerRoutes(app, e)
		
		e.Router.POST("/api/newsletter/subscribe", func(c *core.RequestEvent) error {
			email := strings.TrimSpace(c.Request.FormValue("email"))
			source := strings.TrimSpace(c.Request.FormValue("source"))

			if email == "" {
				return c.HTML(http.StatusBadRequest, `<div class="text-red-400 text-xs mt-2">Email is required.</div>`)
			}
			parsed, err := mail.ParseAddress(email)
			if err != nil || parsed.Address != email {
				return c.HTML(http.StatusBadRequest, `<div class="text-red-400 text-xs mt-2">Please enter a valid email address.</div>`)
			}

			if existing, _ := app.FindFirstRecordByFilter("subscribers", "email = {:email}", dbx.Params{"email": email}); existing != nil {
				return c.HTML(http.StatusOK, `<div class="text-filament text-xs mt-2 font-mono">You are already subscribed. Thanks!</div>`)
			}

			collection, err := app.FindCollectionByNameOrId("subscribers")
			if err != nil {
				app.Logger().Error("Failed to find subscribers collection: " + err.Error())
				return c.HTML(http.StatusInternalServerError, `<div class="text-red-400 text-xs mt-2">Internal server error.</div>`)
			}

			record := core.NewRecord(collection)
			record.Set("email", email)
			record.Set("active", true)
			record.Set("source", source)

			if err := app.Save(record); err != nil {
				app.Logger().Error("Failed to save subscriber: " + err.Error())
				return c.HTML(http.StatusInternalServerError, `<div class="text-red-400 text-xs mt-2">Could not complete subscription. Please try again.</div>`)
			}

			return c.HTML(http.StatusOK, `<div class="text-[#3DDC84] font-medium text-xs mt-2 font-mono">✓ Subscribed. We will keep you updated.</div>`)
		})

		e.Router.POST("/api/contact/submit", func(c *core.RequestEvent) error {
			name := strings.TrimSpace(c.Request.FormValue("name"))
			contactMethod := strings.TrimSpace(c.Request.FormValue("contact_method"))
			subject := strings.TrimSpace(c.Request.FormValue("subject"))
			message := strings.TrimSpace(c.Request.FormValue("message"))

			if name == "" || contactMethod == "" || message == "" {
				return c.HTML(http.StatusBadRequest, `<div class="p-4 bg-red-500/10 border border-red-500 text-red-400 rounded-lg text-xs font-mono">Please fill in all required fields.</div>`)
			}

			collection, err := app.FindCollectionByNameOrId("contact_submissions")
			if err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="p-4 bg-red-500/10 border border-red-500 text-red-400 rounded-lg text-xs font-mono">Internal server error.</div>`)
			}

			record := core.NewRecord(collection)
			record.Set("name", name)
			record.Set("contact_method", contactMethod)
			record.Set("subject", subject)
			record.Set("message", message)
			record.Set("status", "new")

			if err := app.Save(record); err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="p-4 bg-red-500/10 border border-red-500 text-red-400 rounded-lg text-xs font-mono">Failed to save inquiry.</div>`)
			}

			// Fire and forget email notification, logging failures
			go func() {
				if err := mailer.SendContactNotification(name, contactMethod, subject, message); err != nil {
					app.Logger().Error("Failed to send contact notification email: " + err.Error())
				}
			}()

			return c.HTML(http.StatusOK, `<div class="p-6 text-center bg-[#14181F] border border-[#3DDC84]/50 rounded-xl"><h3 class="font-display font-bold text-xl text-[#3DDC84] mb-2">Inquiry Received</h3><p class="text-xs text-[#9DA3AF] leading-relaxed">Thank you, `+html.EscapeString(name)+`. Our team in Karachi will reply via WhatsApp or email shortly.</p></div>`)
		})
		
		// Serve static assets from ui/static
		e.Router.GET("/static/{path...}", func(c *core.RequestEvent) error {
			path := strings.TrimPrefix(c.Request.URL.Path, "/static/")
			http.ServeFile(c.Response, c.Request, filepath.Join("ui", "static", path))
			return nil
		})

		// Serve Service Worker from root
		e.Router.GET("/sw.js", func(c *core.RequestEvent) error {
			http.ServeFile(c.Response, c.Request, filepath.Join("ui", "static", "sw.js"))
			return nil
		})

		// Explicit detail routes
		e.Router.GET("/blog/{slug}", func(c *core.RequestEvent) error {
			return serveContentDetail(app, c, "blog", "blog_detail.html")
		})
		e.Router.GET("/projects/{slug}", func(c *core.RequestEvent) error {
			return serveContentDetail(app, c, "project", "project_detail.html")
		})
		e.Router.GET("/learn/{slug}", func(c *core.RequestEvent) error {
			return serveContentDetail(app, c, "learn", "learn_detail.html")
		})
		e.Router.GET("/products/{slug}", func(c *core.RequestEvent) error {
			return serveProductDetail(app, c)
		})
		e.Router.GET("/printers/{slug}", func(c *core.RequestEvent) error {
			return serveProductDetail(app, c)
		})

		// Catch-all route for static pages and top-level dynamic pages
		e.Router.GET("/{path...}", func(c *core.RequestEvent) error {
			reqPath := c.Request.URL.Path

			// Ignore API and Admin UI requests so PocketBase handles them normally
			if strings.HasPrefix(reqPath, "/api/") || strings.HasPrefix(reqPath, "/_/") {
				return c.Next()
			}

			if reqPath == "/" {
				reqPath = "/home"
			}

			// 1. Try to serve exact path or {path}.html from pb_public
			if info, err := os.Stat(filepath.Join("pb_public", reqPath)); err == nil && !info.IsDir() {
				http.ServeFile(c.Response, c.Request, filepath.Join("pb_public", reqPath))
				return nil
			}
			if info, err := os.Stat(filepath.Join("pb_public", reqPath+".html")); err == nil && !info.IsDir() {
				http.ServeFile(c.Response, c.Request, filepath.Join("pb_public", reqPath+".html"))
				return nil
			}

			// 2. Dynamic page render
			pageName := resolvePageAlias(strings.TrimPrefix(reqPath, "/"))
			siteMap := loadSiteMap(app)
			data := compiler.TemplateData{
				Site: siteMap,
				List: loadPageList(app, pageName),
			}

			htmlOut, err := compiler.RenderTemplate(pageName+".html", data)
			if err == nil {
				return c.HTML(http.StatusOK, htmlOut)
			}

			app.Logger().Error("Template render error for " + pageName + ": " + err.Error())
			return serveNotFound(c, siteMap)
		})

		return e.Next()
	})
}

func resolvePageAlias(name string) string {
	switch name {
	case "printers":
		return "products"
	case "schools-and-universities":
		return "schools"
	case "find-my-printer":
		return "quiz"
	case "fyp", "print-service":
		return "print_service"
	case "parts":
		return "filament"
	case "courses":
		return "workshops"
	default:
		return name
	}
}

func loadSiteMap(app core.App) map[string]any {
	siteSettings, _ := app.FindFirstRecordByFilter("site_settings", "1=1")
	siteMap := make(map[string]any)
	if siteSettings != nil {
		siteMap = siteSettings.PublicExport()
	}
	if _, ok := siteMap["whatsapp_business_number"]; !ok || siteMap["whatsapp_business_number"] == "" {
		siteMap["whatsapp_business_number"] = "923001318112"
	}
	return siteMap
}

func loadPageList(app core.App, pageName string) []map[string]any {
	var listData []map[string]any
	switch pageName {
	case "blog", "projects", "learn":
		cType := pageName
		if pageName == "projects" {
			cType = "project"
		}
		records, _ := app.FindRecordsByFilter("content", "type = {:type} && published = true", "-published_at", 100, 0, dbx.Params{"type": cType})
		for _, r := range records {
			listData = append(listData, r.PublicExport())
		}
	case "products":
		records, _ := app.FindRecordsByFilter("products", "active = true", "-created", 100, 0)
		for _, r := range records {
			listData = append(listData, r.PublicExport())
		}
	}
	return listData
}

func serveContentDetail(app core.App, c *core.RequestEvent, cType, tmplName string) error {
	siteMap := loadSiteMap(app)
	slug := c.Request.PathValue("slug")
	record, _ := app.FindFirstRecordByFilter("content", "type = {:type} && slug = {:slug} && published = true", dbx.Params{"type": cType, "slug": slug})
	if record == nil {
		return serveNotFound(c, siteMap)
	}

	data := compiler.TemplateData{
		Site:   siteMap,
		Record: record.PublicExport(),
	}
	if bodyStr, ok := data.Record["body"].(string); ok {
		rendered, _ := compiler.RenderMarkdown(bodyStr)
		data.Body = compiler.SafeHTML(rendered)
	}

	htmlOut, err := compiler.RenderTemplate(tmplName, data)
	if err != nil {
		app.Logger().Error("Template render error for " + tmplName + ": " + err.Error())
		return serveNotFound(c, siteMap)
	}
	return c.HTML(http.StatusOK, htmlOut)
}

func serveProductDetail(app core.App, c *core.RequestEvent) error {
	siteMap := loadSiteMap(app)
	slug := c.Request.PathValue("slug")
	record, _ := app.FindFirstRecordByFilter("products", "slug = {:slug}", dbx.Params{"slug": slug})
	if record == nil {
		return serveNotFound(c, siteMap)
	}

	recMap := record.PublicExport()
	data := compiler.TemplateData{
		Title:  fmt.Sprintf("%v Price in Pakistan", recMap["name"]),
		Site:   siteMap,
		Record: recMap,
	}

	htmlOut, err := compiler.RenderTemplate("product_detail.html", data)
	if err != nil {
		app.Logger().Error("Template render error for product_detail.html: " + err.Error())
		return serveNotFound(c, siteMap)
	}
	return c.HTML(http.StatusOK, htmlOut)
}

func serveNotFound(c *core.RequestEvent, siteMap map[string]any) error {
	html404, _ := compiler.RenderTemplate("404.html", compiler.TemplateData{Site: siteMap})
	return c.HTML(http.StatusNotFound, html404)
}
