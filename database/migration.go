package database

import (
	"fmt"
	"log"
	"project-with-fiber/model/entity"
)

func RunMigration() {
	err := DB.AutoMigrate(&entity.User{}, &entity.Book{}, &entity.Category{}, &entity.Photo{})
	if err != nil {
		log.Println(err)
	}
	
	fmt.Println("database migrated")
}