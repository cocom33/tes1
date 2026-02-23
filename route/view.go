package route

import (
	"project-with-fiber/handler"

	"github.com/gofiber/fiber/v2"
)


func RouteView(r *fiber.App) {
	
	r.Get("/view", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"title" : "ini title apakah",
		})
	})
	
	r.Get("/category", handler.CategoryHandleGetAll)
}