package populatedb

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func AddProductCategory(req *fiber.Ctx) error {
	var reqBody models.AddProductCategory
	var categoryContent []models.ProductCategory
	if err := req.BodyParser(&reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(&reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	categoryContent = append(categoryContent, reqBody.Categories...)
	addCategory := handlers.DB.Create(&categoryContent)
	if addCategory.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "sorry category cannot be added!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "categories have been added!",
	})
}

func GetCategories(req *fiber.Ctx) error {
	var categories []models.ProductCategory
	queryCategories := handlers.DB.Find(&categories)
	if queryCategories.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no records found!",
		})
	}
	return req.Status(201).JSON(categories)
}
