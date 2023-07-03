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
	getProduct := DB.First(&products, productId)
	if getProduct.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product is either sold or deleted!",
		})
	}
	cartId := GenID()

	updateCartId := DB.Model(&models.User{}).
		Where(&models.User{ID: uint(validation.DecodedToken["id"].(float64))}).
		Update("cart_id", uint(cartId))
	if updateCartId.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product cannot be added to cart__ db error",
		})
	}
	addCart := models.Cart{
		Userid:           uint(validation.DecodedToken["id"].(float64)),
		ProductId:        products.ProductID,
		ProductName:      products.ProductName,
		ProductBrand:     products.ProductBrand,
		ProductCondition: products.ProductCondition,
		ShoeSize:         products.ShoeSize,
		ClothSize:        products.ClothSize,
		Color:            products.Color,
		Price:            products.Price,
	}

	createCart := DB.Create(&addCart)
	if createCart.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "Product cannot be added to cart!",
		})
	}
	var cart []models.Cart
	getCart := DB.Where(&models.Cart{Userid: uint(validation.DecodedToken["id"].(float64))}).
		Find(&cart)
	if getCart.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	var totalPrice float32
	for _, value := range cart {
		totalPrice = totalPrice + value.Price
	}
	response := models.CartResponse{
		CartId:     uint(cartId),
		TotalPrice: totalPrice,
		Items:      cart,
	}
	return req.Status(201).JSON(response)

}

func GetCart(req *fiber.Ctx) error {
	var cart []models.Cart
	var user models.User
	getCart := DB.Preload("Products").Where(&models.Cart{Userid: uint(validation.DecodedToken["id"].(float64))}).
		Find(&cart)
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
	getCartId := DB.Select("cart_id").
		First(&user, uint(validation.DecodedToken["id"].(float64)))
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

func DeleteCartItem(req *fiber.Ctx) error {
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
	deleteCart := DB.Where(&models.Cart{ProductId: productId}).Delete(&models.Cart{})
	if deleteCart.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "the product does not exists!",
		})
	}
	var cart []models.Cart
	var user models.User
	getCart := DB.Preload("Products").Where(&models.Cart{Userid: uint(validation.DecodedToken["id"].(float64))}).
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
		First(&user, uint(validation.DecodedToken["id"].(float64)))
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
