package main

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New()
	db := models.DbConnect()
	handlers.DB = db
	rdb := ConnectRedis()
	handlers.Rds = rdb
	handlers.ConnectSmtp()

	// register routes
	RegisterRoutes(app)

	app.Listen(":3000")
}
