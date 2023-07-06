package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GenID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := rand.Intn(100000)
	return randomInt
}

func AddToCart(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	if len(req.Params("product_id")) == 0 {
		req.Status(400).JSON(fiber.Map{
			"msg": "product_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("product_id"), 10, 32)
	if err != nil {
		return err
	}
	products := new(models.Products)
	productId := uint(num)
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
	addCart := models.Cart{
		Userid:           userId,
		ProductId:        products.ProductID,
		ProductName:      products.ProductName,
		ProductBrand:     products.ProductBrand,
		ProductCondition: products.ProductCondition,
		ShoeSize:         products.ShoeSize,
		ClothSize:        products.ClothSize,
		Color:            products.Color,
		Price:            products.Price,
	}

	//insert new item to cart
	if createCart := DB.Create(&addCart).Error; createCart != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product cannot be added to cart!",
		})
	}
	//Get the list of item in the user's cart
	var cart []models.Cart
	if getCart := DB.Where(&models.Cart{Userid: userId}).Find(&cart).Error; getCart != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	//calculate the total price of items in the user's cart
	var totalPrice float32
	for _, value := range cart {
		totalPrice = totalPrice + value.Price
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
	var cart []models.Cart
	var user models.User
	getCart := DB.Where(&models.Cart{Userid: userId}).Find(&cart)
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
		req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty",
		})
	}
	var totalPrice float32
	cartId := user.CartId
	for _, value := range cart {
		totalPrice = totalPrice + value.Price
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
	deleteCart := DB.Where(&models.Cart{Userid: userId, ProductId: productId}).Delete(&models.Cart{})
	if deleteCart.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "this product does not exists in your cart!",
		})
	}
	var cart []models.Cart
	var user models.User
	getCart := DB.Preload("Products").Where(&models.Cart{Userid: userId}).
		Find(&cart)
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
	getCartId := DB.Select("cart_id").
		First(&user, userId)
	if getCartId.Error != nil {
		req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty",
		})
	}
	var totalPrice float32
	cartId := user.CartId
	for _, value := range cart {
		totalPrice = totalPrice + value.Price
	}

	response := models.CartResponse{
		CartId:     cartId,
		TotalPrice: totalPrice,
		Items:      cart,
	}
	return req.Status(201).JSON(response)
}
