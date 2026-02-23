package handler

import (
	"log"
	"net/http"
	"project-with-fiber/database"
	"project-with-fiber/model/entity"
	"project-with-fiber/model/request"
	"project-with-fiber/ultis"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserHandlerGetAll(ctx *fiber.Ctx) error {
	var users = []entity.User{}
	// userInfo := ctx.Locals("userInfo")
	// log.Println("user info data :: ", userInfo)
	
	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Println(result.Error)
	}
	
	// err := database.DB.Find(&users).Error
	// if err != nil {
	// 	log.Println(err)
	// }
	
	return ctx.JSON(users)
}

func UserHandlerCreate(ctx *fiber.Ctx) error {
	user := new(request.UserCreateRequest)
	validate := validator.New()
	
	err := ctx.BodyParser(user)
	if err != nil {
		log.Println(err)
	}
	
	
	err = validate.Struct(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed",
			"error": strings.Split(err.Error(), "\n"),
		})
	}
	
	users := []request.UserEmailUpdateRequest{}
	err = database.DB.Table("users").Select("email").Scan(&users).Error;
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed",
			"error": strings.Split(err.Error(), "\n"),
		})
	}
	
	for _, v := range users {
		if user.Email == v.Email {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "failed",
				"error": "Email already taken",
			})
		}
	}
	
	pw, err := ultis.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		ctx.Status(500).JSON(fiber.Map{
			"status" : "failed",
			"message" : "internal server error",
		})
	}
	
	newUser := entity.User{
		Name: user.Name,
		Address: user.Address,
		Email: user.Email,
		Password: pw,
		Phone: user.Phone,
	}
	
	err = database.DB.Create(&newUser).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": newUser,
	})
}

func UserHandlerFindById(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	
	var user entity.User
	err := database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "failed",
			"error": "user not found",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": user,
	})
}

func UserHandlerUpdate(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	userRequest := new(request.UserUpdateRequest)
	user := new(entity.User)
	validate := validator.New()

	err := ctx.BodyParser(userRequest)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed",
			"error": "bad request",
		})
	}	
	
	err = database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "failed",
			"error": "user not found",
		})
	}
	
	err = validate.Struct(userRequest)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed",
			"error": strings.Split(err.Error(), "\n"),
		})
	}
	
	if userRequest.Name != "" {
		user.Name = userRequest.Name
	}
	if userRequest.Address != "" {
		user.Address = userRequest.Address
	}
	if userRequest.Phone != "" {
		user.Phone = userRequest.Phone
	}
	
	err = database.DB.Save(user).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed",
			"error": "internal server error",
		})
	}
	
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
		"message": "success",
		"data": user,
	})
}

func UserHandlerUpdateEmail(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	user := new(entity.User)
	userRequest := new(request.UserEmailUpdateRequest)
	validate := validator.New()
	
	err := ctx.BodyParser(userRequest)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed",
			"error": "bad request",
		})
	}
	
	err = validate.Struct(userRequest)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed",
			"error": strings.Split(err.Error(), "\n"),
		})
	}
	
	err = database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "failed",
			"error": "not found",
		})
	}
	
	users := []request.UserEmailUpdateRequest{}
	err = database.DB.Table("users").Where("email != ?", user.Email).Select("email").Scan(&users).Error;
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed",
			"error": strings.Split(err.Error(), "\n"),
		})
	}
	
	for _, v := range users {
		if userRequest.Email == v.Email {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "failed",
				"error": "Email already taken",
			})
		}
	}
	
	if userRequest.Email != "" {
		user.Email = userRequest.Email
	}
	
	err = database.DB.Save(user).Error
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed",
			"error": "internal server error",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": user,
	})
}

func UserHandlerDelete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	var user entity.User
	
	// check available user
	err := database.DB.First(&user, "id=?", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	
	err = database.DB.Debug().Delete(&user).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"message": "user was deleted",
	})
}