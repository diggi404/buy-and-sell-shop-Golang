package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func Signup(req *fiber.Ctx) error {
	reqBody := new(models.SignupSchema)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)

	}
	checkUser := DB.Where(&models.User{Email: reqBody.Email}).First(&models.User{})
	if checkUser.Error != nil {
		hash, _ := validation.HashPassword(reqBody.Password)
		results := DB.Create(&models.User{Name: reqBody.Name, Email: reqBody.Email, Password: hash})
		if results.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "signup failed!",
			})
		}
		return req.Status(201).JSON(fiber.Map{
			"msg": "Welcome to Golang!",
		})
	}
	return req.Status(400).JSON(fiber.Map{
		"msg": "user already exists!",
	})
}
