package handler

import (
	"log"
	"net/http"
	"project-with-fiber/database"
	"project-with-fiber/model/entity"
	"project-with-fiber/model/request"
	"project-with-fiber/ultis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func BookHandlerGetAll(ctx *fiber.Ctx) error {
	// Get all books from database
	book := new([]entity.Book)
	
	res := database.DB.Table("books").Find(book)
	if res.Error != nil {
		log.Println(res.Error)
	}
	
	return ctx.JSON(book)
}

func BookHandlerGetById(ctx *fiber.Ctx) error {
	// Get book by id from database
	book := new(entity.Book)
	id := ctx.Params("id")
	
	res := database.DB.Table("books").First(&book, "id = ?", id)
	if res.Error != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"status" : "failed",
			"message" : "book not found",
		})
	}
	
	return ctx.JSON(book)
}

func BookHandlerCreate(ctx *fiber.Ctx) error {
	book := new(request.BookRequest)
	if err := ctx.BodyParser(book); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failed",
			"message": err.Error(),
		})
	}
	
	validate := validator.New()
	err := validate.Struct(book)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failed",
			"message": err.Error(),
		})
	}
	
	// handle file
	fileName, err := ultis.HandleSingleFile(ctx, "cover", "/images")
	if err != nil {
		log.Println(err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	newBook := entity.Book{
		Title : book.Title,
		Cover : fileName,
		Author : book.Author,
	}
	
	errCreateBook := database.DB.Create(&newBook).Error
	if errCreateBook != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status": "failed",
			"message": errCreateBook.Error(),
		})
	}
	
	return ctx.JSON(fiber.Map{
		"status": "success",
		"data": newBook,
	})
}