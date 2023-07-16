package main

import (
	handlers "Users/diggi/Documents/Go_tutorials/handlers"
	populatedb "Users/diggi/Documents/Go_tutorials/handlers/populateDB"
	usersettings "Users/diggi/Documents/Go_tutorials/handlers/userSettings"
	validation "Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {

	app.Post("/auth/login", handlers.Login)
	app.Post("/signup", handlers.Signup)
	app.Post("/verify/email/:link", handlers.ConfirmEmail)
	app.Get("/user/profile", validation.Authenticator, handlers.UserProfile)
	app.Put("/user/email", validation.Authenticator, usersettings.UpdateEmail)
	app.Post("/user/email/code", validation.Authenticator, usersettings.VerifyEmailOtp)
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
	app.Delete("/user/payments/:card_id", validation.Authenticator, handlers.DeleteCrediCard)
	app.Get("/user/payments_methods", validation.Authenticator, handlers.GetPaymentMethods)
	app.Put("/user/payments_methods/:card_id", validation.Authenticator, handlers.MakeCardDefault)
	app.Post("/user/checkout", validation.Authenticator, handlers.Checkout)
	app.Get("/user/orders", validation.Authenticator, handlers.GetUserOrders)
	app.Get("/user/orders/in-progress", validation.Authenticator, handlers.GetInProgressItems)
	app.Put("/user/orders/tracking/:item_id", validation.Authenticator, handlers.FixTrackingNumbers)
}
