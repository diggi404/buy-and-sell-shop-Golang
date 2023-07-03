package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PostItem(req *fiber.Ctx) error {
	reqBody := new(models.AddProduct)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	productsContent := models.Products{
		UserID:           uint(validation.DecodedToken["id"].(float64)),
		ProductName:      reqBody.ProductName,
		Categoryid:       reqBody.Categoryid,
		ProductBrand:     reqBody.ProductBrand,
		ProductCondition: reqBody.ProductCondition,
		ShoeSize:         reqBody.ShoeSize,
		ClothSize:        reqBody.ClothSize,
		Color:            reqBody.Color,
		Price:            reqBody.Price,
	}

	addProduct := DB.Create(&productsContent)
	if addProduct.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "cannot add product details!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "Item has been successfully posted!",
	})
}

func GetUserProducts(req *fiber.Ctx) error {
	var products []models.Products
	query := models.Products{UserID: uint(validation.DecodedToken["id"].(float64))}
	getProducts := DB.Where(&query).Find(&products)
	if getProducts.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getProducts.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no products found!",
		})
	}
	return req.Status(201).JSON(products)
}

func DeleteProduct(req *fiber.Ctx) error {
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
	deleteItem := DB.
		Where(&models.Products{UserID: uint(validation.DecodedToken["id"].(float64)), ProductID: productId}).
		Delete(&models.Products{})
	if deleteItem.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "the product does not exists!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "product has been deleted successfully!",
	})

}

// func UpdateProduct(req *fiber.Ctx) error {
// 	if len(req.Params("product_id")) == 0 {
// 		return req.Status(400).JSON(fiber.Map{
// 			"msg": "product_id is required!",
// 		})
// 	}
// 	num, err := strconv.ParseUint(req.Params("product_id"), 10, 32)
// 	if err != nil {
// 		return err
// 	}
// 	product_id := uint(num)
// 	reqBody := new(models.AddProduct)
// 	if err := req.BodyParser(reqBody); err != nil {
// 		return err
// 	}
// 	errors := validation.ValidateStruct(*reqBody)
// 	if errors != nil {
// 		return req.Status(400).JSON(errors)
// 	}
// 	productsContent := models.Products{
// 		UserID:           uint(validation.DecodedToken["id"].(float64)),
// 		ProductName:      reqBody.ProductName,
// 		Categoryid:       reqBody.Categoryid,
// 		ProductBrand:     reqBody.ProductBrand,
// 		ProductCondition: reqBody.ProductCondition,
// 		ShoeSize:         reqBody.ShoeSize,
// 		ClothSize:        reqBody.ClothSize,
// 		Color:            reqBody.Color,
// 		Price:            reqBody.Price,
// 	}
// 	updateItem := DB.Save()
// }
