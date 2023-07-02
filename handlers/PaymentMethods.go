package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddCreditCard(req *fiber.Ctx) error {
	var creditCard models.CreditCard
	reqBody := new(models.CreditCard)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}

	creditCard = models.CreditCard{
		User_id:    uint(validation.DecodedToken["id"].(float64)),
		CardNumber: reqBody.CardNumber,
		CardMonth:  reqBody.CardMonth,
		CardYear:   reqBody.CardYear,
	}
	addCard := DB.Create(&creditCard)
	if addCard.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error adding credit card!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "credit card has been saved!",
	})
}

func DeleteCrediCard(req *fiber.Ctx) error {
	if len(req.Params("card_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "card_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("card_id"), 10, 32)
	if err != nil {
		return err
	}
	cardId := uint(num)
	deleteCard := DB.Delete(&models.CreditCard{}, cardId)
	if deleteCard.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error deleting credit card!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "card has been deleted!",
	})
}
