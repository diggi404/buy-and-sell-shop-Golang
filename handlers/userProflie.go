package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"log"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func UserProfile(req *fiber.Ctx) error {
	dbResponse := new(models.User)
	getUser := DB.First(&dbResponse, validation.DecodedToken["id"])
	if getUser.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "user does not exists",
		})
	} else {
		mailer := gomail.NewMessage()
		mailer.SetAddressHeader("From", "karianfavreau9@gmail.com", "Buy Sell")
		mailer.SetAddressHeader("To", "dbackson1@gmail.com", "")
		mailer.SetHeader("Subject", "Password Reset")
		mailer.SetBody("text/plain", "Here is your code for resetting your password")
		if err := Smtp.DialAndSend(mailer); err != nil {
			log.Fatal(err)
		}
		return req.Status(201).JSON(fiber.Map{
			"name":  dbResponse.Name,
			"email": dbResponse.Email,
		})
	}
}
