package main

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	populatedb "Users/diggi/Documents/Go_tutorials/handlers/populateDB"
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
	// db.Debug().AutoMigrate(&models.User{}, &models.AddressBook{}, &models.CategoryOptions{}, &models.Products{}, &models.ProductCategory{}, &models.Cart{})
	// db.Debug().AutoMigrate(&models.Products{}, &models.Cart{})
	// db.Debug().AutoMigrate(&models.User{}, &models.Cart{})
	// db.Debug().AutoMigrate(&models.CreditCard{}, &models.MobileMoney{}, &models.BillingAddress{}, &models.TotalCart{})
	// db.Debug().AutoMigrate(&models.Products{}, &models.ProductCategory{}, &models.CategoryOptions{}, &models.Cart{})
	// db.Debug().AutoMigrate(&models.User{}, &models.AddressBook{}, &models.CreditCard{})
	// db.Debug().AutoMigrate(&models.CreditCard{})
	handlers.DB = db
	app.Post("/auth/login", handlers.Login)
	app.Post("/signup", handlers.Signup)
	app.Get("/user/profile", validation.Authenticator, handlers.UserProfile)
	app.Post("/user/create/address", validation.Authenticator, handlers.CreteAddressBook)
	app.Get("/user/address/", validation.Authenticator, handlers.GetAddressBook)
	app.Put("/user/update/address/:address_id", validation.Authenticator, handlers.UpdateAddressBook)
	app.Post("/create/category", validation.Authenticator, populatedb.AddProductCategory)
	app.Get("/categories", validation.Authenticator, populatedb.GetCategories)
	app.Post("/create/options", validation.Authenticator, populatedb.AddCategoryOptions)
	app.Get("/options/:category_id", validation.Authenticator, populatedb.GetCategoryOptions)
	app.Post("/user/create/item", validation.Authenticator, handlers.PostItem)
	app.Get("/user/items", validation.Authenticator, handlers.GetUserProducts)
	app.Get("/products/all", validation.Authenticator, handlers.GetAllProducts)
	app.Delete("/user/item/:product_id", validation.Authenticator, handlers.DeleteProduct)
	app.Post("/user/create/cart/:product_id", validation.Authenticator, handlers.AddToCart)
	app.Get("/user/cart", validation.Authenticator, handlers.GetCart)
	app.Delete("/user/cart/:product_id", validation.Authenticator, handlers.DeleteCartItem)
	app.Post("/user/create/credit_card", validation.Authenticator, handlers.AddCreditCard)
	// app.Get("/user/credit_cards", validation.Authenticator, handlers.GetCreditCards)
	app.Delete("/user/payments/:card_id", validation.Authenticator, handlers.DeleteCrediCard)
	app.Get("/user/payments_methods", validation.Authenticator, handlers.PaymentMethods)

	app.Listen(":3000")
}
