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
	checkEmtpyCart := DB.Preload("Products").Where(&models.Cart{Userid: userId}).Find(&cart)
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
		AddressID:     reqBody.AddressId,
		PaymentMethod: paymentInfo,
		PaidTotal:     totalCart.TotalPrice,
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
	var (
		m        []map[string]interface{}
		prods    []map[string]interface{}
		prodIds  []models.Products
		products []models.Products
	)
	for _, c := range cart {
		products = append(products, c.Products)
	}
	cartByte, _ := json.Marshal(products)
	json.Unmarshal(cartByte, &m)
	for _, value := range m {
		v := make(map[string]interface{})
		value["order_id"] = order.OrderId
		v["product_id"] = value["product_id"]

		prods = append(prods, v)
	}
	prodByte, _ := json.Marshal(prods)
	json.Unmarshal(prodByte, &prodIds)

	newCartByte, _ := json.Marshal(m)

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
	var (
		tempShipmentMap []map[string]interface{}
		shipmentMap     []map[string]interface{}
		shipment        []models.Shipment
	)
	itemsByte, _ := json.Marshal(items)
	json.Unmarshal(itemsByte, &tempShipmentMap)
	for _, item := range tempShipmentMap {
		temp := make(map[string]interface{})
		temp["item_id"] = item["item_id"]
		shipmentMap = append(shipmentMap, temp)
	}
	shipmentByte, _ := json.Marshal(shipmentMap)
	json.Unmarshal(shipmentByte, &shipment)
	DB.Create(&shipment)
	deletSoldProduct := DB.Delete(&prodIds)
	if deletSoldProduct.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred. Please try again!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "order has been confirmed!",
	})

}
