package handlers

import (
	"Users/diggi/Documents/Go_tutorials/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ConfirmEmail(req *fiber.Ctx) error {
	if len(req.Params("link")) == 0 {
		return req.Status(400).JSON(fiber.Map{
			"msg": "required paramet is missing!",
		})
	}
	var link models.EmailVerify
	if checkLinkErr := DB.Where(&models.EmailVerify{Link: req.Params("link")}).
		First(&link).Error; checkLinkErr != nil {
		return req.Status(401).JSON(fiber.Map{
			"msg": "your link is invalid!",
		})
	}
	if link.ExpiresAt.Unix() < time.Now().Unix() {
		DB.Where(&models.EmailVerify{Link: req.Params("link")}).Delete(&models.EmailVerify{})
		return req.Status(401).JSON(fiber.Map{
			"msg": "your link is expired!",
		})
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		deleteLink := tx.Where(&models.EmailVerify{Link: req.Params("link")}).
			Delete(&models.EmailVerify{})
		if deleteLink.RowsAffected == 0 {
			return deleteLink.Error
		}
		updateVerify := tx.Model(&models.User{}).Where(&models.User{ID: link.UserId}).Update("email_verified", true)
		if updateVerify.RowsAffected == 0 {
			return updateVerify.Error
		}
		return nil
	})
	if err != nil {
		return req.Status(400).JSON(fiber.Map{
			"msg": "an error occurred!",
		})
	}
	return req.Status(201).JSON(fiber.Map{
		"msg": "your email has been verified!",
	})
}
