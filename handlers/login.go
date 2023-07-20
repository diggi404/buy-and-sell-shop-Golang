package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Rds *redis.Client

func Login(req *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
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
	if !authUser.EmailVerified {
		return req.Status(401).JSON(fiber.Map{
			"msg": "Your email is not yet verified! Kindly check your mail inbox to proceed!",
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
	token, err := validation.GenerateJwt(os.Getenv("SECRET_KEY"), payload)
	if err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "error generating token!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"access_token": token,
	})

}
