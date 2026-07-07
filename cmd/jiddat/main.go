package main

import (
	"log"
	"mime"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"jiddat3d/internal/hooks"
	"jiddat3d/internal/routes"
	_ "jiddat3d/pb_migrations"
)

func main() {
	// Explicitly register MIME types in case the host OS (like Alpine) is missing /etc/mime.types
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".svg", "image/svg+xml")

	app := pocketbase.New()

	// Register migration commands if not building in strict mode
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	routes.RegisterPublicRoutes(app)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		hooks.RegisterContentHooks(app)
		
		// Register placeholder healthz route
		e.Router.GET("/healthz", func(e *core.RequestEvent) error {
			return e.String(200, "OK")
		})

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
