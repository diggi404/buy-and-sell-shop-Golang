package populatedb

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func AddCategoryOptions(req *fiber.Ctx) error {
	var reqBody models.ValidateCategoryOptions
	var optionsContent []models.CategoryOptions
	if err := req.BodyParser(&reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(&reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	optionsContent = append(optionsContent, reqBody.Options...)
	addOptions := handlers.DB.Create(&optionsContent)
	if addOptions.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "options cannot be added!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "options have added successfully!",
	})
}
