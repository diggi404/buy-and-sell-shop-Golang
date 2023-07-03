package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Signup(req *fiber.Ctx) error {
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
		results := DB.Create(&models.User{Name: reqBody.Name, Email: reqBody.Email, Password: hash})
		if results.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "signup failed!",
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
	return req.Status(400).JSON(fiber.Map{
		"msg": "user already exists!",
	})
}
