package route

import (
	"project-with-fiber/handler"
	"project-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	
	r.Static("/public", "./public/assets")
	r.Static("/images", "./public/images")

	r.Post("/login", handler.LoginHandler)
	r.Post("/register", handler.RegisterHandler)
	r.Post("/logout", handler.LogoutHandler)
	r.Get("/profile", handler.LoginHandlerProfile, middleware.UserMiddleware)
	
	r.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic : error")
	})

	user := r.Group("/user", middleware.UserMiddleware)
	user.Get("/", handler.UserHandlerGetAll)
	user.Post("/", handler.UserHandlerCreate)
	user.Get("/:id", handler.UserHandlerFindById)
	user.Put("/:id", handler.UserHandlerUpdate)
	user.Put("/:id/update-email", handler.UserHandlerUpdateEmail)
	user.Delete("/:id", handler.UserHandlerDelete)
	
	book := r.Group("/book", middleware.UserMiddleware)
	book.Post("/", handler.BookHandlerCreate)
	book.Get("/", handler.BookHandlerGetAll)
	book.Get("/:id", handler.BookHandlerGetById)
	
	photo := r.Group("/photo", middleware.UserMiddleware)
	photo.Post("/", handler.PhotoHandlerCreate)
	photo.Delete("/:categoryId", middleware.CheckBody, handler.PhotoHandlerDelete)
	
	category := r.Group("/category", middleware.UserMiddleware)
	category.Post("/", handler.CategoryHandleCreate)
	category.Get("/", handler.CategoryHandleGetAll)
}