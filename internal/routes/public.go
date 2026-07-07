package routes

import (
	"strings"
	"net/http"

	"jiddat3d/internal/mailer"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterPublicRoutes(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		
		e.Router.POST("/api/newsletter/subscribe", func(c *core.RequestEvent) error {
			email := strings.TrimSpace(c.Request.FormValue("email"))
			source := strings.TrimSpace(c.Request.FormValue("source"))
			
			if email == "" {
				return c.HTML(http.StatusBadRequest, `<div class="text-accent-terracotta text-sm mt-2">Email is required.</div>`)
			}

			// We rely on PocketBase's unique constraint on email.
			// If it fails, they are likely already subscribed.

			// Create record
			collection, err := app.FindCollectionByNameOrId("subscribers")
			if err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="text-accent-terracotta text-sm mt-2">Internal server error.</div>`)
			}

			record := core.NewRecord(collection)
			record.Set("email", email)
			record.Set("active", true)
			record.Set("source", source)

			if err := app.Save(record); err != nil {
				// Assume duplicate email
				return c.HTML(http.StatusOK, `<div class="text-accent-gold text-sm mt-2">You are already subscribed. Thanks!</div>`)
			}

			return c.HTML(http.StatusOK, `<div class="text-primary font-medium text-sm mt-2">You're in. We'll keep you posted.</div>`)
		})

		e.Router.POST("/api/contact/submit", func(c *core.RequestEvent) error {
			name := strings.TrimSpace(c.Request.FormValue("name"))
			contactMethod := strings.TrimSpace(c.Request.FormValue("contact_method"))
			subject := strings.TrimSpace(c.Request.FormValue("subject"))
			message := strings.TrimSpace(c.Request.FormValue("message"))

			if name == "" || contactMethod == "" || message == "" {
				return c.HTML(http.StatusBadRequest, `<div class="p-4 bg-accent-terracotta/10 border border-accent-terracotta text-accent-terracotta rounded-md">Please fill in all required fields.</div>`)
			}

			collection, err := app.FindCollectionByNameOrId("contact_submissions")
			if err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="p-4 bg-accent-terracotta/10 border border-accent-terracotta text-accent-terracotta rounded-md">Internal server error.</div>`)
			}

			record := core.NewRecord(collection)
			record.Set("name", name)
			record.Set("contact_method", contactMethod)
			record.Set("subject", subject)
			record.Set("message", message)
			record.Set("status", "new")

			if err := app.Save(record); err != nil {
				return c.HTML(http.StatusInternalServerError, `<div class="p-4 bg-accent-terracotta/10 border border-accent-terracotta text-accent-terracotta rounded-md">Failed to send message.</div>`)
			}

			// Fire and forget email notification
			go mailer.SendContactNotification(name, contactMethod, subject, message)

			return c.HTML(http.StatusOK, `<div class="p-8 text-center"><h3 class="font-serif text-2xl text-primary mb-2">Message Sent</h3><p class="text-ink-muted">We will get back to you shortly.</p></div>`)
		})
		
		return e.Next()
	})
}
