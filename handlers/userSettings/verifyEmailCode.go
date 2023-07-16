package usersettings

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifyEmailOtp(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	reqBody := new(models.VerifyCode)
	if err := req.BodyParser(reqBody); err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "invalid request Body",
		})
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)

	}
	var otp models.OtpEmail
	getOtp := handlers.DB.Where(&models.OtpEmail{UserId: userId, Value: reqBody.Code}).First(&otp)
	if getOtp.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "invalid code!",
		})
	}
	if otp.ExpiresAt.Unix() < time.Now().Unix() {
		handlers.DB.Delete(&models.OtpEmail{}, userId)
		return req.Status(400).JSON(fiber.Map{
			"msg": "your code has expired!",
		})
	}
	handlers.DB.Delete(&models.OtpEmail{}, userId)
	return req.Status(201).JSON(fiber.Map{
		"msg": "code has been verified!",
	})
}
