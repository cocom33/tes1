package handler

import (
	"fmt"
	"log"
	"net/http"
	"project-with-fiber/database"
	"project-with-fiber/model/entity"
	"project-with-fiber/model/request"
	"project-with-fiber/ultis"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)
	validate := validator.New()
	
	err := ctx.BodyParser(loginRequest)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"status": "failed",
			"message": err.Error(),
		})
	}
	
	// validasi request
	err = validate.Struct(loginRequest)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status": "failed",
			"message": err.Error(),
		})
	}
	
	// check user
	var user entity.User
	err = database.DB.Where("email = ?", loginRequest.Email).First(&user).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "failed",
			"message": "wrong credentials",
		})
	}
	
	// check password
	isValid := ultis.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "failed",
			"message": "wrong credentials",
		})
	}
	
	// generate jwt
	claims := jwt.MapClaims{}
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["address"] = user.Address
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	
	token, err := ultis.GenerateToken(&claims)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "failed",
			"message": "internal server error",
		})
	}
	
	session, err := database.Store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create session",
		})
	}
	
	session.Set("userId", user.ID)
	session.Set("userToken", token)
	err = session.Save()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
	}
	
	// tt := session.Get("userId")
	// ut := session.Get("userToken")
	// fmt.Println(tt)
	// fmt.Println(ut)
	// fmt.Println("Session Cookie:", ctx.Cookies("fiber_sess"))
	
	return ctx.JSON(fiber.Map{
		"status" : "success",
		"token" : token,
	})
}

func LogoutHandler(ctx *fiber.Ctx) error {
	session, err := database.Store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not get session",
		})
	}
	
	userToken := session.Get("userToken")
	if userToken == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	err = session.Destroy()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not logout",
		}) 
	}

	return ctx.JSON(fiber.Map{
		"status": "success",
		"message" : "berhasil logout",
	})
}

func LoginHandlerProfile(ctx *fiber.Ctx) error {
	sesi, err := database.Store.Get(ctx) 
	if err != nil {
		return err
	}
	fmt.Println(ctx.Cookies("fiber_sess"))
	
	user := new(entity.User)
	err = database.DB.First(user, "id = ?", sesi.Get("userId")).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "failed",
			"message": "Unauthorized",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"status": "success",
		"profile": fiber.Map{
			"data": user,
			"token": sesi.Get("userToken"),
		},
	})
}

func RegisterHandler(ctx *fiber.Ctx) error {
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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