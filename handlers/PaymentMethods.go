package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	creditcard "github.com/durango/go-credit-card"
	"github.com/gofiber/fiber/v2"
)

func AddCreditCard(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var creditCard models.CreditCard
	reqBody := new(models.CreditCard)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	ccNum := strconv.FormatUint(uint64(reqBody.CardNumber), 10)
	card := creditcard.Card{
		Number: ccNum,
	}
	if err := card.ValidateNumber(); !err {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your credit card number is invalid!",
		})
	}

	creditCard = models.CreditCard{
		User_ID:    userId,
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

func GetCreditCards(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var creditcards []models.CreditCard
	getCards := DB.Where(&models.CreditCard{User_ID: userId}).Find(&creditcards)
	if getCards.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getCards.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no saved credit cards found!",
		})
	}

	return req.Status(201).JSON(creditcards)
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
	if deleteCard.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error deleting credit card!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "card has been deleted!",
	})
}

func PaymentMethods(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var users []models.User
	getPaymentMethods := DB.Preload("CreditCards").Preload("Momo").First(&users, userId)
	if getPaymentMethods.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no user found!",
		})
	}
	return req.Status(201).JSON(users)
}
