package main

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	populatedb "Users/diggi/Documents/Go_tutorials/handlers/populateDB"
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
	app := fiber.New()
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USERNAME") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database!")
	}
	db.Debug().AutoMigrate(&models.CategoryOptions{})
	handlers.DB = db
	app.Post("/auth/login", handlers.Login)
	app.Post("/signup", handlers.Signup)
	app.Get("/user/profile", validation.Authenticator, handlers.UserProfile)
	app.Post("/user/create/address", validation.Authenticator, handlers.CreteAddressBook)
	app.Put("/user/update/address/:address_id", validation.Authenticator, handlers.UpdateAddressBook)
	app.Post("/create/category", validation.Authenticator, populatedb.AddProductCategory)
	app.Get("/categories", validation.Authenticator, populatedb.GetCategories)
	app.Post("/create/options", validation.Authenticator, populatedb.AddCategoryOptions)
	app.Get("/options/:category_id", validation.Authenticator, populatedb.GetCategoryOptions)
	app.Post("/user/create/item", validation.Authenticator, handlers.PostItem)

	app.Listen(":3000")
}
