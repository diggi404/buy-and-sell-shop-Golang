package main

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	app := fiber.New()
	dsn := "host=" + os.Getenv("DB_HOST") + "user=" + os.Getenv("DB_USERNAME") + "password=" + os.Getenv("DB_PASSWORD") + "dbname=" + os.Getenv("DB_NAME") + "port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database!")
	}
	db.Debug().AutoMigrate(&models.AddressBook{})
	handlers.DB = db
	app.Post("/auth/login", handlers.Login)
	app.Post("/signup", handlers.Signup)
	app.Get("/user/profile", validation.Authenticator, handlers.UserProfile)
	app.Post("/user/create/address", validation.Authenticator, handlers.CreteAddressBook)

	app.Listen(":3000")
}
