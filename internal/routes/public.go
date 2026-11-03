package routes

import (
	"fmt"
	"net/http"
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
				return c.HTML(http.StatusBadRequest, `<div class="text-[#FF6A1A] text-xs mt-2">Email is required.</div>`)
			}

			collection, err := app.FindCollectionByNameOrId("subscribers")
			if err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="text-[#FF6A1A] text-xs mt-2">Internal server error.</div>`)
			}

			record := core.NewRecord(collection)
			record.Set("email", email)
			record.Set("active", true)
			record.Set("source", source)

			if err := app.Save(record); err != nil {
				return c.HTML(http.StatusOK, `<div class="text-[#FFB400] text-xs mt-2 font-mono">You are already subscribed. Thanks!</div>`)
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

			// Fire and forget email notification
			go mailer.SendContactNotification(name, contactMethod, subject, message)

			return c.HTML(http.StatusOK, `<div class="p-6 text-center bg-[#14181F] border border-[#3DDC84]/50 rounded-xl"><h3 class="font-display font-bold text-xl text-[#3DDC84] mb-2">Inquiry Received</h3><p class="text-xs text-[#9DA3AF] leading-relaxed">Thank you, `+name+`. Our team in Karachi will reply via WhatsApp or email shortly.</p></div>`)
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

		// Catch-all route to serve static/cached templates, mimicking Caddy's try_files
		e.Router.GET("/{path...}", func(c *core.RequestEvent) error {
			reqPath := c.Request.URL.Path
			
			// Ignore API and Admin UI requests so PocketBase handles them normally
			if strings.HasPrefix(reqPath, "/api/") || strings.HasPrefix(reqPath, "/_/") {
				return c.Next()
			}
			
			if reqPath == "/" {
				reqPath = "/home"
			}
			
			// 1. Try to serve exact path from pb_public
			if info, err := os.Stat(filepath.Join("pb_public", reqPath)); err == nil && !info.IsDir() {
				http.ServeFile(c.Response, c.Request, filepath.Join("pb_public", reqPath))
				return nil
			}
			
			// 2. Try {path}.html
			if info, err := os.Stat(filepath.Join("pb_public", reqPath+".html")); err == nil && !info.IsDir() {
				http.ServeFile(c.Response, c.Request, filepath.Join("pb_public", reqPath+".html"))
				return nil
			}
			
			// 3. Try /cached{path}.html
			if info, err := os.Stat(filepath.Join("pb_public", "cached", reqPath+".html")); err == nil && !info.IsDir() {
				http.ServeFile(c.Response, c.Request, filepath.Join("pb_public", "cached", reqPath+".html"))
				return nil
			}
			
			// 4. Try dynamic render
			rawPath := strings.TrimPrefix(reqPath, "/")
			pageName := rawPath
			
			// Path aliases
			switch pageName {
			case "printers":
				pageName = "products"
			case "schools-and-universities":
				pageName = "schools"
			case "find-my-printer":
				pageName = "quiz"
			case "fyp":
				pageName = "print_service"
			case "print-service":
				pageName = "print_service"
			case "parts":
				pageName = "filament"
			case "courses":
				pageName = "workshops"
			}
			
			siteSettings, _ := app.FindFirstRecordByFilter("site_settings", "1=1")
			siteMap := make(map[string]any)
			if siteSettings != nil {
				siteMap = siteSettings.PublicExport()
			}
			if _, ok := siteMap["whatsapp_business_number"]; !ok || siteMap["whatsapp_business_number"] == "" {
				siteMap["whatsapp_business_number"] = "923001318112"
			}
			
			var listData []map[string]any
			if pageName == "blog" || pageName == "projects" || pageName == "learn" {
				cType := pageName
				if pageName == "projects" {
					cType = "project"
				}
				records, _ := app.FindRecordsByFilter("content", "type = {:type} && published = true", "-published_at", 100, 0, dbx.Params{"type": cType})
				for _, r := range records {
					listData = append(listData, r.PublicExport())
				}
			} else if pageName == "products" {
				records, _ := app.FindRecordsByFilter("products", "active = true", "-created", 100, 0)
				for _, r := range records {
					listData = append(listData, r.PublicExport())
				}
			}
			
			data := compiler.TemplateData{
				Site: siteMap,
				List: listData,
			}
			
			// Handle Detail Pages dynamically
			if strings.HasPrefix(reqPath, "/blog/") || strings.HasPrefix(reqPath, "/projects/") || strings.HasPrefix(reqPath, "/learn/") {
				parts := strings.Split(strings.TrimPrefix(reqPath, "/"), "/")
				if len(parts) == 2 {
					cType := parts[0]
					if cType == "projects" {
						cType = "project"
					}
					slug := parts[1]
					record, _ := app.FindFirstRecordByFilter("content", "type = {:type} && slug = {:slug} && published = true", dbx.Params{"type": cType, "slug": slug})
					if record != nil {
						data.Record = record.PublicExport()
						if bodyStr, ok := data.Record["body"].(string); ok {
							rendered, _ := compiler.RenderMarkdown(bodyStr)
							data.Body = compiler.SafeHTML(rendered)
						}
						
						tmplName := parts[0] + "_detail.html"
						if parts[0] == "projects" {
							tmplName = "project_detail.html"
						}
						
						html, err := compiler.RenderTemplate(tmplName, data)
						if err == nil {
							return c.HTML(http.StatusOK, html)
						}
					}
				}
			}

			// Handle Product / Printer Detail Pages
			if strings.HasPrefix(reqPath, "/products/") || strings.HasPrefix(reqPath, "/printers/") {
				parts := strings.Split(strings.TrimPrefix(reqPath, "/"), "/")
				if len(parts) == 2 {
					slug := parts[1]
					record, _ := app.FindFirstRecordByFilter("products", "slug = {:slug}", dbx.Params{"slug": slug})
					if record != nil {
						data.Record = record.PublicExport()
					} else {
						// Fallback catalog map for specialized Jiddat machines
						fallbackCatalog := map[string]map[string]any{
							"jiddat-apex-corexy": {
								"name":        "Jiddat Apex CoreXY",
								"tagline":     "Industrial enclosed high-speed FDM/FFM workhorse running custom Klipper firmware. Built for engineering plastics (ABS, Nylon, PC).",
								"price":       "245,000",
								"currency":    "PKR",
								"availability": "built-to-order",
								"featured":    true,
							},
							"jiddat-mark-1": {
								"name":        "Jiddat Mark-1 Desktop",
								"tagline":     "Precision direct-drive FDM/FFM printer with CNC-machined rigid chassis and all-metal hotend for makers and engineers.",
								"price":       "145,000",
								"currency":    "PKR",
								"availability": "in-stock",
								"featured":    true,
							},
							"jiddat-edu-maker": {
								"name":        "Jiddat Edu-Maker V2",
								"tagline":     "Classroom-safe, rugged educational FDM/FFM printer designed for Pakistani schools, colleges, and university STEAM labs.",
								"price":       "98,000",
								"currency":    "PKR",
								"availability": "in-stock",
								"featured":    true,
							},
							"jiddat-l500-workstation": {
								"name":        "Jiddat L-500 Large-Format",
								"tagline":     "Massive 500x500x500mm build chamber for industrial tooling, full-scale automotive parts, and final year project prototypes.",
								"price":       "385,000",
								"currency":    "PKR",
								"availability": "built-to-order",
								"featured":    true,
							},
							"jiddat-corexy": {
								"name":        "Jiddat CoreXY Maker Edition",
								"tagline":     "100% open-source RepRap DIY maker assembly kit with local precision hardware, open firmware, and full CAD schematics.",
								"price":       "78,000",
								"currency":    "PKR",
								"availability": "in-stock",
								"featured":    false,
							},
							"jiddat-precision-resin": {
								"name":        "Jiddat Precision Resin (MSLA)",
								"tagline":     "Ultra-high resolution stereolithography system for jewelry and dental casting — Currently in Engineering R&D in Karachi.",
								"price":       "135,000",
								"currency":    "PKR",
								"availability": "in-development",
								"featured":    false,
							},
						}

						if mockItem, exists := fallbackCatalog[slug]; exists {
							data.Record = mockItem
						}
					}

					if data.Record != nil {
						data.Title = fmt.Sprintf("%v Price in Pakistan", data.Record["name"])
						html, err := compiler.RenderTemplate("product_detail.html", data)
						if err == nil {
							return c.HTML(http.StatusOK, html)
						}
					}
				}
			}
			
			html, err := compiler.RenderTemplate(pageName+".html", data)
			if err == nil {
				return c.HTML(http.StatusOK, html)
			}
			
			// Log the error!
			app.Logger().Error("Template render error for " + pageName + ": " + err.Error())
			
			// Fallback 404
			html404, _ := compiler.RenderTemplate("404.html", data)
			return c.HTML(http.StatusNotFound, html404)
		})
		
		return e.Next()
	})
}
