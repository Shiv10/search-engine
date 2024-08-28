package routes

import (
	"search-engine/db"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func SetRoutes(app *fiber.App) {
	app.Get("/", AuthMiddleware, DashboardHandler)

	app.Post("/", AuthMiddleware, DashboardPostHandler)

	app.Get("/login", LoginHandler)

	app.Post("/login", LoginPostHandler)

	app.Get("/create", func(c *fiber.Ctx) error {
		u := &db.User{}
		u.CreateAdmin()
		return c.SendString("created")
	})
}
