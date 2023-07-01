package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Login(req *fiber.Ctx) error {
	authUser := new(models.User)
	reqBody := new(models.LoginSchema)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)

	}
	checkEmail := DB.Where(&models.User{Email: reqBody.Email}).First(&authUser)
	if checkEmail.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "invalid email or password!",
		})
	}
	if checkHash := validation.CheckPasswordHash(reqBody.Password, authUser.Password); !checkHash {
		return req.Status(401).JSON(fiber.Map{
			"msg": "invalid email or password!",
		})
	}
	payload := &jwt.MapClaims{
		"id":    authUser.ID,
		"email": authUser.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := validation.GenerateJwt("fadfadsfasf", payload)
	if err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error generating token!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"access_token": token,
	})

}
