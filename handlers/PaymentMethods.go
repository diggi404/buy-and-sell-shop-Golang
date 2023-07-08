package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	"github.com/ggwhite/go-masker"

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

	ccNum := strconv.FormatUint(uint64(reqBody.CardNumber), 10)
	ccMon := strconv.FormatUint(uint64(reqBody.CardMonth), 10)
	ccYear := strconv.FormatUint(uint64(reqBody.CardYear), 10)

	// execute when the user doesn't select an address from the addressbook
	if len(reqBody.AddressID) == 0 {
		addressReqBody := new(models.CreditCardManual)
		if err := req.BodyParser(addressReqBody); err != nil {
			return err
		}
		errors := validation.ValidateStruct(*addressReqBody)
		if errors != nil {
			return req.Status(400).JSON(errors)
		}
		card := creditcard.Card{
			Number: ccNum,
			Month:  ccMon,
			Year:   ccYear,
		}
		if err := card.ValidateNumber(); !err {
			return req.Status(400).JSON(fiber.Map{
				"msg": "your credit card number is invalid!",
			})
		}
		err := card.Method()
		if err != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "your credit card is invalid!",
			})
		}
		rawLastFour, _ := card.LastFour()
		cardType := card.Company.Short
		lastFour, _ := strconv.ParseUint(rawLastFour, 10, 32)
		maskedCC := masker.CreditCard(ccNum)
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
		addressId := addresses.AddressId
		creditCard = models.CreditCard{
			User_ID:      userId,
			AddressID:    addressId,
			CardNumber:   reqBody.CardNumber,
			CardMonth:    reqBody.CardMonth,
			CardYear:     reqBody.CardYear,
			CardType:     cardType,
			LastFour:     uint(lastFour),
			MaskedNumber: maskedCC,
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
	card := creditcard.Card{
		Number: ccNum,
		Month:  ccMon,
		Year:   ccYear,
	}
	if err := card.ValidateNumber(); !err {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your credit card number is invalid!",
		})
	}
	cardError := card.Method()
	if cardError != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your credit card is invalid!",
		})
	}
	rawLastFour, _ := card.LastFour()
	cardType := card.Company.Short
	lastFour, _ := strconv.ParseUint(rawLastFour, 10, 32)
	maskedCC := masker.CreditCard(ccNum)
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
		User_ID:      userId,
		AddressID:    uint(addressId),
		CardNumber:   reqBody.CardNumber,
		CardMonth:    reqBody.CardMonth,
		CardYear:     reqBody.CardYear,
		LastFour:     uint(lastFour),
		CardType:     cardType,
		MaskedNumber: maskedCC,
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
	if deleteCard.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error deleting credit card!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "card has been deleted!",
	})
}

func GetPaymentMethods(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var users []models.User
	getPaymentMethods := DB.Preload("CreditCards").
		Preload("CreditCards.Address").
		Preload("Momo").First(&users, userId)
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

func MakeCardDefault(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	if len(req.Params("card_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "required parameter is missing!",
		})
	}
	card_id, err := strconv.ParseUint(req.Params("card_id"), 10, 32)
	if err != nil {
		return err
	}

	checkCard := DB.Where(&models.CreditCard{
		User_ID: userId,
		CardId:  uint(card_id),
	}).
		First(&models.CreditCard{})
	if checkCard.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no such card found!",
		})
	}
	makeDefault := DB.Model(&models.User{}).
		Where(&models.User{ID: uint(userId)}).
		Update("default_payment_method", card_id)
	if makeDefault.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error making card default!",
		})
	}
	if makeDefault.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no such card found!",
		})
	}
	cleanUpDefault := DB.Model(&models.CreditCard{}).
		Where(&models.CreditCard{IsDefault: true}).
		Update("is_default", false)
	if cleanUpDefault.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error making default!",
		})
	}
	setDefault := DB.Model(&models.CreditCard{}).
		Where(&models.CreditCard{CardId: uint(card_id)}).
		Update("is_default", true)
	if setDefault.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error making card default!",
		})
	}
	if setDefault.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no such card found!aaa",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "default payment registered!",
	})
}
