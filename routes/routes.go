package routes

import (
	"search-engine/views"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type loginform struct {
	Email string `form:"email"`
	Password string `form:"password"`
}

type settingsform struct {
	Amount int `form:"amount"`
	SearchOn bool `form:"searchOn"`
	AddNew bool `form:"addNew"`
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func SetRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, views.Home())
	})

	app.Post("/", func(c *fiber.Ctx) error {
		input := settingsform{}
		if err := c.BodyParser(&input); err != nil {
			c.SendString("<h2>Error: Something went wrong</h2>")
		}
		return c.SendStatus(200)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return render(c, views.Login())
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		input := loginform{}
		if err := c.BodyParser(&input); err != nil {
			c.SendString("<h2>Error: Something went wrong</h2>")
		}
		return c.SendStatus(200)
	})
}
