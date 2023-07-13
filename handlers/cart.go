package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func GenID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := rand.Intn(100000)
	return randomInt
}

func AddToCart(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	if len(req.Params("product_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "product_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("product_id"), 10, 32)
	if err != nil {
		return err
	}
	products, productId := new(models.Products), uint(num)
	//check if product exists before adding to cart
	if getProduct := DB.First(&products, productId).Error; getProduct != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product is either sold or deleted!",
		})
	}
	//Check if item was posted by the same user
	checkOwnership := DB.Where(&models.Products{UserID: userId, ProductID: productId}).
		Find(&models.Products{})
	if checkOwnership.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkOwnership.RowsAffected > 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "you already own this item!",
		})
	}
	//Check if user has already added that item to his cart
	checkCartItemExists := DB.Where(&models.Cart{ProductId: productId}).Find(&models.Cart{})
	if checkCartItemExists.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkCartItemExists.RowsAffected > 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "this item is already in your cart!",
		})
	}
	//first delete cartId to prevent intefering with updating cartId in user table(foreign key errors)
	if delCartId := DB.Where(&models.TotalCart{User_Id: userId}).Delete(&models.TotalCart{}).Error; delCartId != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	//update cart id in the user's table
	cartId := GenID()
	if updateCartId := DB.Model(&models.User{}).
		Where(&models.User{ID: userId}).
		Update("cart_id", uint(cartId)).Error; updateCartId != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product cannot be added to cart!",
		})
	}
	likes := products.ItemLikes + 1
	DB.Model(&models.Products{}).Where(&models.Products{ProductID: productId}).Update("likes", likes)
	addCart := models.Cart{
		Userid:    userId,
		ProductId: products.ProductID,
	}

	//insert new item to cart
	if createCart := DB.Create(&addCart).Error; createCart != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product cannot be added to cart!",
		})
	}
	//Get the list of item in the user's cart
	var cart []models.Cart
	if getCart := DB.Preload("Products").Where(&models.Cart{Userid: userId}).Find(&cart).Error; getCart != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	//calculate the total price of items in the user's cart
	var totalPrice float32
	for _, value := range cart {
		totalPrice += value.Products.Price
	}
	// Save total price and cart id to total_carts_table
	if saveCartInfo := DB.Save(&models.TotalCart{
		CartID:     uint(cartId),
		User_Id:    userId,
		TotalPrice: totalPrice,
	}).Error; saveCartInfo != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	response := models.CartResponse{
		CartId:     uint(cartId),
		TotalPrice: totalPrice,
		Items:      cart,
	}
	return req.Status(201).JSON(response)

}

func GetCart(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var (
		cart []models.Cart
		user models.User
	)
	getCart := DB.Preload("Products").Where(&models.Cart{Userid: userId}).Find(&cart)
	if getCart.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getCart.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty!",
		})
	}
	if getCartId := DB.Select("cart_id").First(&user, userId).Error; getCartId != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty",
		})
	}

	var totalPrice float32
	cartId := user.CartId
	for _, value := range cart {
		totalPrice += value.Products.Price
	}
	response := models.CartResponse{
		CartId:     cartId,
		TotalPrice: totalPrice,
		Items:      cart,
	}
	return req.Status(201).JSON(response)
}

func DeleteCartItem(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	if len(req.Params("product_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "product_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("product_id"), 10, 32)
	if err != nil {
		return err
	}
	productId := uint(num)
	var delCart models.Cart

	// delete the item in the cart
	deleteItem := DB.Clauses(clause.Returning{}).
		Where(&models.Cart{Userid: userId, ProductId: productId}).
		Delete(&delCart)
	if deleteItem.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "this product does not exists in your cart!",
		})
	}

	// update likes count
	likes := delCart.Products.ItemLikes - 1
	DB.Model(&models.Products{}).
		Where(&models.Products{ProductID: productId}).
		Update("likes", likes)
	var (
		cart []models.Cart
		user models.User
	)
	getCart := DB.Preload("Products").Where(&models.Cart{Userid: userId}).Find(&cart)
	if getCart.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getCart.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your cart is now empty!",
		})
	}
	getCartId := DB.Select("cart_id").First(&user, userId)
	if getCartId.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty",
		})
	}
	var totalPrice float32
	cartId := user.CartId
	for _, value := range cart {
		totalPrice = totalPrice + value.Products.Price
	}

	response := models.CartResponse{
		CartId:     cartId,
		TotalPrice: totalPrice,
		Items:      cart,
	}
	return req.Status(201).JSON(response)
}
