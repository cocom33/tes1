package database

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Store *session.Store
	DB *gorm.DB
)

func DatabaseInit() {
	Store = session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   false, // Cookie hanya dikirim melalui HTTPS
		CookieSameSite: "Lax", // Atur berdasarkan kebutuhan
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:fiber_sess",
	})

	var err error
	
	database := "root@tcp(127.0.0.1:3306)/go-project?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(database), &gorm.Config{})
	
	if err != nil {
		panic("can't connect to database")
	}
	
	fmt.Println("Connected to database")
}