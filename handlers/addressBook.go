package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func CreteAddressBook(req *fiber.Ctx) error {
	reqBody := new(models.CreateAddressBook)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	checkCount := DB.Where("user_id = ?", validation.DecodedToken["id"]).Find(&models.AddressBook{})
	if checkCount.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	} else if checkCount.RowsAffected >= 3 {
		return req.Status(401).JSON(fiber.Map{
			"msg": "you have reached your addressBook limit!",
		})
	} else if checkCount.RowsAffected == 0 {
		addressContent := &models.AddressBook{UserId: uint(validation.DecodedToken["id"].(float64)), FirstName: reqBody.FirstName, LastName: reqBody.LastName, Address1: reqBody.Address1, Address2: reqBody.Address2, City: reqBody.City, State: reqBody.State, ZipCode: reqBody.ZipCode}
		registerAdressBook := DB.Create(&addressContent)
		if registerAdressBook.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "sorry address cannot be added!",
			})
		}
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "address has been added successfully!",
	})
}
