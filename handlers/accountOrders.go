package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"strconv"

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

func GetInProgressItems(req *fiber.Ctx) error {
	userId := validation.DecodedToken["id"].(float64)
	var items []models.PurchasedItems
	getItems := DB.Where(&models.PurchasedItems{SellerId: uint(userId), OrderStatus: "processing"}).Preload("Shipment").Find(&items)
	if getItems.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getItems.RowsAffected == 0 {
		return req.Status(201).JSON(fiber.Map{
			"msg": "no orders found!",
		})
	}
	return req.Status(201).JSON(items)
}

func FixTrackingNumbers(req *fiber.Ctx) error {
	reqBody := new(models.Shipment)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	if len(req.Params("item_id")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "required parameter is missing!",
		})
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	itemId, _ := strconv.ParseUint(req.Params("item_id"), 10, 32)
	tracking := models.Shipment{
		TrackingNumber: reqBody.TrackingNumber,
		Carrier:        reqBody.Carrier,
	}
	updateTracking := DB.Model(&models.Shipment{}).Where(&models.Shipment{ItemId: uint(itemId)}).Updates(tracking)
	if updateTracking.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred. Please try again!",
		})
	}
	DB.Model(&models.PurchasedItems{}).
		Where(&models.PurchasedItems{ItemID: uint(itemId)}).
		Updates(map[string]interface{}{"order_status": "shipped", "can_cancel": false})

	return req.Status(201).JSON(fiber.Map{
		"msg": "the tracking number has been updated!",
	})
}
