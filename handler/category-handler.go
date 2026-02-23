package handler

import (
	"net/http"
	"project-with-fiber/database"
	"project-with-fiber/model/entity"
	"project-with-fiber/model/request"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CategoryHandleCreate(ctx *fiber.Ctx) error {
	category := new(request.CategoryRequest)
	if err := ctx.BodyParser(category); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	validate := validator.New()
	err := validate.Struct(category)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	newCategory := entity.Category{
		Name: category.Name,
	}
	
	err = database.DB.Create(&newCategory).Error
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err,
		})
	}
	
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status" : "success",
		"data" : newCategory,
	})
}

func CategoryHandleGetAll(ctx *fiber.Ctx) error {
	category := new([]entity.Category)
	
	err := database.DB.Preload("Photos").Find(&category).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status" : "failed",
			"message" : err,
		})
	}
	
	return ctx.JSON(fiber.Map{
		"status" : "success",
		"data" : category,
	})
}