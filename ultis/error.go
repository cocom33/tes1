package ultis

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Cek apakah error yang terjadi adalah fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		// Kembalikan respon error dengan status code dan pesan error yang sesuai
		return ctx.Status(e.Code).JSON(fiber.Map{
			"error": e.Message,
		})
	}

	// Jika error bukan tipe fiber.Error, tangani dengan respon umum
	log.Println("Internal Server Error:", err)
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal Server Error",
	})
}
