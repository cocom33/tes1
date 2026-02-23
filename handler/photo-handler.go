package handler

import (
	"fmt"
	"log"
	"net/http"
	"project-with-fiber/database"
	"project-with-fiber/model/entity"
	"project-with-fiber/model/request"
	"project-with-fiber/ultis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func PhotoHandlerCreate(ctx *fiber.Ctx) error {
	photo := new(request.PhotoRequest)
	if err := ctx.BodyParser(photo); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	validate := validator.New()
	err := validate.Struct(photo)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	names, err := ultis.HandleMultipleFile(ctx, "photo", "photo")
	if err != nil {
		log.Println(err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error(),
		})
	}
	
	fmt.Println(names)
	fmt.Println(photo.CategoryId)

	var data []entity.Photo
	for _, v := range names {
		newPhoto := &entity.Photo{
			Image : v,
			CategoryId: photo.CategoryId,
		}
		
		data = append(data, *newPhoto)
	}
	
	err = database.DB.Create(&data).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status" : "failed",
			"message" : err,
		})
	}
	
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status" : "success",
		"message" : "photo created",
	})
}

func PhotoHandlerDelete(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")
	photo := []entity.Photo{}
	
	if categoryId == "" {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"status" : "failed",
			"message" : "photo not found",
		})
	}

	err := database.DB.Where("category_id = ?", categoryId).Find(&photo).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status" : "failed",
			"message" : err.Error,
		})
	}
	
	for _, v := range photo {
		err = ultis.HandleRemoveFile(v.Image)
		if err != nil {
			// return err
			// fmt.Println(err)
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"status" : "failed",
				"message" : "File not found.",
			})
		}
	}
	
	err = database.DB.Delete(&photo).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status" : "failed",
			"message" : "internal server error",
		})
	}
	
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status" : "success",
		"message" : "photo deleted",
	})
}