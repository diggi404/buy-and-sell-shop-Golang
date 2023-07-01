package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"Users/diggi/Documents/Go_tutorials/validation"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var addressBook []models.AddressBook

func CreteAddressBook(req *fiber.Ctx) error {
	reqBody := new(models.CreateAddressBook)
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	checkCount := DB.Where(
		&models.AddressBook{UserId: uint(validation.DecodedToken["id"].(float64))}).
		Find(&addressBook)
	if checkCount.Error != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if checkCount.RowsAffected >= 3 {
		return req.Status(401).JSON(fiber.Map{
			"msg": "you have reached your addressBook limit!",
		})
	} else {
		addressContent := &models.AddressBook{
			UserId:    uint(validation.DecodedToken["id"].(float64)),
			FirstName: reqBody.FirstName,
			LastName:  reqBody.LastName,
			Address1:  reqBody.Address1,
			City:      reqBody.City,
			State:     reqBody.State,
			ZipCode:   reqBody.ZipCode,
		}
		registerAdressBook := DB.Create(&addressContent)
		if registerAdressBook.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "sorry address cannot be added!",
			})
		}
		return req.Status(201).JSON(fiber.Map{
			"msg": "address has been added successfully!",
		})
	}

}

func UpdateAddressBook(req *fiber.Ctx) error {
	reqBody := new(models.UpdateAddressBook)
	if len(req.Params("address_id")) == 0 {
		return req.Status(401).JSON(fiber.Map{
			"msg": "address_id is required!",
		})
	}
	if err := req.BodyParser(reqBody); err != nil {
		return err
	}
	errors := validation.ValidateStruct(*reqBody)
	if errors != nil {
		return req.Status(400).JSON(errors)
	}
	num, err := strconv.ParseUint(req.Params("address_id"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	addressId := uint(num)
	checkAddress := DB.First(&addressBook, addressId)
	if checkAddress.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "address cannot be found!",
		})
	} else {
		addressContent := &models.AddressBook{
			AddressId: addressId,
			UserId:    uint(validation.DecodedToken["id"].(float64)),
			FirstName: reqBody.FirstName,
			LastName:  reqBody.LastName,
			Address1:  reqBody.Address1,
			City:      reqBody.City,
			State:     reqBody.State,
			ZipCode:   reqBody.ZipCode}
		updateAddress := DB.Model(&addressBook).Updates(addressContent)
		if updateAddress.Error != nil {
			return req.Status(400).JSON(fiber.Map{
				"msg": "error updating address!",
			})
		}
		return req.Status(201).JSON(fiber.Map{
			"msg": "your address has been updated successfully!",
		})
	}
}

func GetAddressBook(req *fiber.Ctx) error {
	var addressList []models.AddressBook
	getAddress := DB.Where(
		&models.AddressBook{UserId: uint(validation.DecodedToken["id"].(float64))}).
		Find(&addressList)
	if getAddress.Error != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	if getAddress.RowsAffected == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "no records found!",
		})
	}
	return req.Status(201).JSON(addressList)
}
