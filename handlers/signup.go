package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	uuid "github.com/uuid6/uuid6go-proto"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm/clause"
)

func Signup(req *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
	authUser := new(models.User)
	reqBody := new(models.SignupSchema)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)

	}
	checkUser := DB.Where(&models.User{Email: reqBody.Email}).First(&models.User{})
	if checkUser.Error != nil {
		hash, _ := validation.HashPassword(reqBody.Password)
		authUser = &models.User{Name: reqBody.Name, Email: reqBody.Email, Password: hash}
		results := DB.Clauses(clause.Returning{}).Create(&authUser)
		if results.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "signup failed!",
			})
		}
		var gen uuid.UUIDv7Generator
		gen.SubsecondPrecisionLength = 20
		link := gen.Next().ToString()
		mailer := gomail.NewMessage()
		mailer.SetAddressHeader("From", "karianfavreau9@gmail.com", "Tradex")
		mailer.SetAddressHeader("To", authUser.Email, "")
		mailer.SetHeader("Subject", "Confirm Your Email Address")
		mailBody := "click on this link below to confirm your email address\nhttp://localhost:3000/verify/email/" + link + "\nLink is valid for only 24 hours"
		mailer.SetBody("text/plain", mailBody)
		if err := Smtp.DialAndSend(mailer); err != nil {
			return err
		}
		DB.Create(&models.EmailVerify{
			UserId:    authUser.ID,
			Link:      link,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		})
		return req.Status(201).JSON(fiber.Map{
			"msg": "Check your inbox or junk for a link to confirm your email address",
		})
	}
	return req.Status(400).JSON(fiber.Map{
		"msg": "user already exists. Please login!",
	})
}
