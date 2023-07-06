package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"fmt"
	"strconv"

	creditcard "github.com/durango/go-credit-card"
	"github.com/gofiber/fiber/v2"
)

func AddCreditCard(req *fiber.Ctx) error {
	// initialize some variables
	userId := uint(validation.DecodedToken["id"].(float64))
	var addresses models.AddressBook
	var creditCard models.CreditCard
	reqBody := new(models.AddCreditCard)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	// validate request body
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}

	// execute when the user doesn't select an address from the addressbook
	if len(reqBody.AddressID) == 0 {
		newReqBody := new(models.CreateAddressBook)
		errors := validation.ValidateStruct(*newReqBody)
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
		checkCard := DB.Where(&models.CreditCard{User_ID: userId, CardNumber: reqBody.CardNumber}).
			Find(&models.CreditCard{})
		if checkCard.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "an error occurred!",
			})
		}
		if checkCard.RowsAffected != 0 {
			return req.Status(400).JSON(fiber.Map{
				"msg": "this card has already been saved to your account!",
			})
		}
		newAddress := models.AddressBook{
			FirstName: reqBody.Address.FirstName,
			LastName:  reqBody.Address.LastName,
			Address1:  reqBody.Address.Address1,
			City:      reqBody.Address.City,
			State:     reqBody.Address.State,
			ZipCode:   reqBody.Address.ZipCode,
		}
		createAddress := DB.Create(&newAddress)
		if createAddress.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "an error occurred!",
			})
		}
		fmt.Println(addresses)
		addressId := addresses.AddressId
		creditCard = models.CreditCard{
			User_ID:    userId,
			AddressID:  addressId,
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

	//execute when the user chooses an address from the addressbook
	addressId, err := strconv.ParseUint(reqBody.AddressID, 10, 32)
	if err != nil {
		return err
	}
	checkCard := DB.Where(&models.CreditCard{User_ID: userId, CardNumber: reqBody.CardNumber}).
		Find(&models.CreditCard{})
	if checkCard.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkCard.RowsAffected != 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "this card has already been saved to your account!",
		})
	}
	creditCard = models.CreditCard{
		User_ID:    userId,
		AddressID:  uint(addressId),
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
	getPaymentMethods := DB.Preload("CreditCards").Preload("CreditCards.Address").Preload("Momo").First(&users, userId)
	if getPaymentMethods.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getPaymentMethods.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no user found!",
		})
	}
	return req.Status(201).JSON(users)
}
