package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func UserProfile(req *fiber.Ctx) error {
	dbResponse := new(models.User)
	getUser := DB.First(&dbResponse, validation.DecodedToken["id"])
	if getUser.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "user does not exists",
		})
	} else {
		return req.Status(201).JSON(fiber.Map{
			"name":  dbResponse.Name,
			"email": dbResponse.Email,
		})
	}
}
