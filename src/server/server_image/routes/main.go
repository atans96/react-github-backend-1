package routes

import (
	"github.com/gofiber/fiber/v2"
	"server_go/src/service"
)

type Routes struct {
	*service.Server
}

func (app *Routes) StartRoutes() {
	// Routes
	app.Options("/images_from_markdown", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/images_from_markdown", ImagesFromMarkdown)
}
