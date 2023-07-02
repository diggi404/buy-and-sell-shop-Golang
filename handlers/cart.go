package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
	return req.Status(201).JSON(fiber.Map{
		"msg": "Product has been successfully added to cart!",
	})

}

func GetCart(req *fiber.Ctx) error {
	var cart []models.Cart
	getCart := DB.Where(&models.Cart{Userid: uint(validation.DecodedToken["id"].(float64))}).Find(&cart)
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
	return req.Status(201).JSON(cart)
}

func DeleteCart(req *fiber.Ctx) error {
	if len(req.Params("product_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "cart_id is required!",
		})
	}
	num, err := strconv.ParseUint(req.Params("product_id"), 10, 32)
	if err != nil {
		return err
	}
	productId := uint(num)
	deleteCart := DB.Where(&models.Cart{ProductId: productId}).Delete(&models.Cart{})
	if deleteCart.Error != nil {
		req.Status(400).JSON(fiber.Map{
			"msg": "error deleting item!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "item has been deleted!",
	})
}
