package populatedb

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddCategoryOptions(req *fiber.Ctx) error {
	var reqBody models.AddCategoryOptions
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

func GetCategoryOptions(req *fiber.Ctx) error {
	if len(req.Params("category_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "category_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("category_id"), 10, 32)
	if err != nil {
		return err
	}
	category_id := uint(num)
	var options []models.CategoryOptions
	query := models.CategoryOptions{CategoryID: category_id}
	getOptions := handlers.DB.Where(&query).Find(&options)
	if getOptions.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getOptions.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no records found!",
		})
	}
	return req.Status(201).JSON(options)
}
