package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Checkout(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var (
		cart       []models.Cart
		totalCart  models.TotalCart
		creditcard models.CreditCard
		address    models.AddressBook
	)
	reqBody := new(models.CreditCardCheckout)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	checkEmtpyCart := DB.Where(&models.Cart{Userid: userId}).Find(&cart)
	if checkEmtpyCart.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkEmtpyCart.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "your cart is empty!",
		})
	}
	confirmTotalPrice := DB.Where(&models.TotalCart{CartID: reqBody.CartId}).First(&totalCart)
	if confirmTotalPrice.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "sorry your cart is empty!",
		})
	}
	if totalCart.TotalPrice != reqBody.TotalPrice {
		return req.Status(400).JSON(fiber.Map{
			"msg": "the total price provided is not valid!",
		})
	}
	checkCard := DB.Where(&models.CreditCard{CardId: reqBody.CardId}).First(&creditcard)
	if checkCard.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no credit card found!",
		})
	}
	checkAddress := DB.Where(&models.AddressBook{AddressId: reqBody.AddressId}).First(&address)
	if checkAddress.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "address does not exists!",
		})
	}
	lastFour := strconv.FormatUint(uint64(creditcard.LastFour), 10)
	paymentInfo := creditcard.CardType + " ending in " + lastFour
	order := models.Orders{
		UserId:        userId,
		PaymentId:     creditcard.CardId,
		AddressId:     address.AddressId,
		PaymentMethod: paymentInfo,
		PaidTotal:     totalCart.TotalPrice,
		Status:        "processing",
		CanCancel:     true,
	}
	keepOrder := DB.Create(&order)
	if keepOrder.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if keepOrder.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred. Please try again!",
		})
	}

	// marshal cart struct into json byte
	cartByte, _ := json.Marshal(cart)

	//unmarsh the byte into a map to easily add the order_id
	var m []map[string]interface{}
	json.Unmarshal(cartByte, &m)
	for _, value := range m {
		value["order_id"] = order.OrderId
	}
	//then marsh the map into a new json byte
	newCartByte, _ := json.Marshal(m)

	// unmarshal the byte into the purchasedItem struct
	var items []models.PurchasedItems
	json.Unmarshal(newCartByte, &items)

	keepPurchasedItems := DB.Create(&items)
	if keepPurchasedItems.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if keepPurchasedItems.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred. Please try again!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "order has been confirmed!",
	})

}

func GetUserOrders(req *fiber.Ctx) error {
	userId := uint(validation.DecodedToken["id"].(float64))
	var orders []models.Orders
	checkOrders := DB.Preload("PurchasedItems").
		Preload("AddressBook").
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
