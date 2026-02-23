package ultis

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
)

const defaultPath = "./public/"

func storeFile(ctx *fiber.Ctx, file *multipart.FileHeader, dirPath string) (string, error) {
	// Generate a random name
	name := ksuid.New().String() + filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s/%s", dirPath, name)
	
	// create if folder doesn't exist
	if err := os.MkdirAll(fmt.Sprintf("./public/%s/", dirPath), os.ModePerm); err != nil {
		return "", errors.New("fail to create directory")
	}
	
	err := ctx.SaveFile(file, fmt.Sprintf("./public/%s", fileName))
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("fail to store file into directory ./public/" + dirPath)
	}
	
	return fileName, nil
}

func contains(str string) bool {
	typeImage := []string{"image/jpg", "image/png", "image/jpeg", "image/gif"}
	
	for _, v := range typeImage {
		if v == str {
			return true
		}
	}
	return false
}

func HandleSingleFile(ctx *fiber.Ctx, nameField string, dirPath string) (string, error) {
	file, errFile := ctx.FormFile(nameField)
	if errFile != nil {
		log.Println("Error file : ", errFile)
	}
	
	if file != nil {
		typeFile := file.Header.Get("Content-Type")
		if !contains(typeFile) {
			return "", errors.New("this type file is not allowed")
		}
		
		fileName, err := storeFile(ctx, file, dirPath)
		if err != nil {
			return "", err
		}
		
		return fileName, nil
	}

	log.Println("nothing file to upload")
	return "", nil
}

func HandleMultipleFile(ctx *fiber.Ctx, nameField string, dirPath string) ([]string, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println("Error form : ", err)
		return nil, err
	}
	
	files := form.File[nameField]
	var filenames []string 

	// typeFile := files.Header.Get("Content-Type")
	// if !contains(typeFile) {
	// 	return nil, errors.New("this type file is not allowed")
	// }
	
	for _, file := range files {
		name, err := storeFile(ctx, file, dirPath)
		if err != nil {
			return nil, err
		}
		
		filenames = append(filenames, name)
	}
	
	return filenames, nil
}

func HandleRemoveFile(filename string, pathFile ...string) error {
	if len(pathFile) > 0 {
		err := os.Remove(pathFile[0] + filename)
		if err != nil {
			log.Println("Error remove file : ", err)
			return err
		}
	} else {
		err := os.Remove(defaultPath + filename)
		if err != nil {
			log.Println("Error remove file : ", err)
			return err
		}
	}
	
	return nil
}