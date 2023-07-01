package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func PostItem(req *fiber.Ctx) error {
	reqBody := new(models.AddProduct)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	productsContent := models.Products{
		UserID:           reqBody.UserId,
		ProductName:      reqBody.ProductName,
		Categoryid:       reqBody.Categoryid,
		ProductBrand:     reqBody.ProductBrand,
		ProductCondition: reqBody.ProductCondition,
		Price:            reqBody.Price,
	}

	addProduct := DB.Create(&productsContent)
	if addProduct.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "cannot add product details!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "Item has been successfully posted!",
	})
}
