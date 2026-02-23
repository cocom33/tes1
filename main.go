package main

import (
	"log"
	"project-with-fiber/database"
	"project-with-fiber/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var engine = html.New("./template", ".html")

func main() {
	// database init
	database.DatabaseInit()
	database.RunMigration()
	
	app := fiber.New(fiber.Config{
		Views:  engine,
	})
	
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				log.Printf("Recovered from panic: %v", r)

				// Return a custom error response
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Something went wrong!",
				})
			}
		}()
		return c.Next() // Proceed to the next middleware/handler
	})

	// route init
	route.RouteInit(app)
	route.RouteView(app)
	
	app.Listen(":8080")
}