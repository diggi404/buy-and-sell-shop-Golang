package usersettings

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func GenOtp() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := 100000 + rand.Intn(900000)
	return randomInt
}

func UpdateEmail(req *fiber.Ctx) error {
	var user models.User
	userId := uint(validation.DecodedToken["id"].(float64))
	getEmail := handlers.DB.Where(&models.User{ID: userId}).First(&user)
	if getEmail.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "user does not exists!",
		})
	}
	mailer := gomail.NewMessage()
	mailer.SetAddressHeader("From", "karianfavreau9@gmail.com", "Tradex")
	mailer.SetAddressHeader("To", user.Email, "")
	mailer.SetHeader("Subject", "Confirm Email Update")
	otpCode := GenOtp()
	mailer.SetBody("text/plain", strconv.Itoa(otpCode)+"\nUse this code to verify your email update. Code expires in 10 min")
	if err := handlers.Smtp.DialAndSend(mailer); err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred. Mail cannot be sent!",
		})
	}
	otpCreds := models.OtpEmail{
		UserId:    userId,
		Value:     otpCode,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(time.Minute * 2),
	}
	err := handlers.DB.Transaction(func(tx *gorm.DB) error {
		checkOtp := handlers.DB.Where(&models.OtpEmail{UserId: userId}).First(&models.OtpEmail{})
		if checkOtp.Error != nil {
			handlers.DB.Create(&otpCreds)
			return nil
		}
		updateOtp := handlers.DB.Where(&models.OtpEmail{UserId: userId}).
			Update("value", otpCode)
		if updateOtp.RowsAffected == 0 {
			return updateOtp.Error
		}
		return nil
	})
	if err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred.Please try again later",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "Code has been sent to your email. Please verify to proceed!",
	})
}

// func ResendCode(req *fiber.Ctx) error {
// 	userId := uint(validation.DecodedToken["id"].(float64))
// 	findUser := handlers.DB.Where(&models.EmailVerify{UserId: userId}).First(&models.EmailVerify{})
// 	if findUser.Error != nil {
// 		return req.Status(400).JSON(fiber.Map{
// 			"msg": "user has not requested for email update!",
// 		})
// 	}
// 	mailer := gomail.NewMessage()
// 	mailer.SetAddressHeader("From", "karianfavreau9@gmail.com", "Tradex")
// 	mailer.SetAddressHeader("To", user.Email, "")
// 	mailer.SetHeader("Subject", "Confirm Email Update")
// 	otpCode := GenOtp()
// 	mailer.SetBody("text/plain", strconv.Itoa(otpCode)+"\nUse this code to verify your email update. Code expires in 10 min")
// 	if err := handlers.Smtp.DialAndSend(mailer); err != nil {
// 		return req.Status(400).JSON(fiber.Map{
// 			"msg": "an error occurred. Mail cannot be sent!",
// 		})
// 	}
// }
