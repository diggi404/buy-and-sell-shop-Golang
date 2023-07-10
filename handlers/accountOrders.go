package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"

	"github.com/gofiber/fiber/v2"
)

func GetUserOrders(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var orders []models.Orders
	checkOrders := DB.Preload("PurchasedItems").Preload("AddressBook").
		Preload("PurchasedItems.Shipment").
		Preload("PurchasedItems.Seller").
		Where(&models.Orders{UserId: userId}).
		Find(&orders)
	if checkOrders.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkOrders.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your order history is empty. Kindly make a purchase!",
		})
	}
	return req.Status(201).JSON(orders)
}
