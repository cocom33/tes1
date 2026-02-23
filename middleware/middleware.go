package middleware

import (
	"fmt"
	"log"
	"project-with-fiber/ultis"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CheckBody(ctx *fiber.Ctx) error {
	log.Println("Request Body:", string(ctx.Body()))
	return ctx.Next()	
}

func UserMiddleware(ctx *fiber.Ctx) error {
	token := ctx.Get("x-token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status" : "failed",
			"message": "unauthenticated",
		})
	}

	claims, err := ultis.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status" : "failed",
			"message": "invalid token",
		})
	}
	
	// role := claims["role"].(string)
	// if role != "admin" {
	// 	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"status" : "failed",
	// 		"message": "forbidden access",
	// 	})
	// }
	
	info := ctx.Locals("userInfo", claims).(jwt.MapClaims)
	fmt.Println(info)
	

	return ctx.Next()
}



// func PermissionCreate() 