package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(255);not null;unique"`
	Password  string `gorm:"type:varchar(255);not null"`
	createdAt time.Time
	updatedAt time.Time
}

type AddressBook struct {
	AddressId uint   `gorm:"primaryKey;autoIncrement; column:address_id"`
	UserId    uint   `gorm:"column:user_id"`
	User      User   `gorm:"foreignKey:UserId"`
	FirstName string `gorm:"varchar(100);not null; column:fname"`
	LastName  string `gorm:"varchar(100);not null; column:lname"`
	Address1  string `gorm:"varchar(255);not null"`
	City      string `gorm:"varchar(100);not null"`
	State     string `gorm:"varchar(100);notn null"`
	ZipCode   string `gorm:"varchar(50);not null"`
	createdAt time.Time
	updatedAt time.Time
}

type ProductCategory struct {
	CategoryId   uint   `gorm:"primaryKey;autoIncrement; column:category_id"`
	CategoryName string `json:"category_name" gorm:"varchar(255);not null; column:category_name"`
	createdAt    time.Time
	updatedAt    time.Time
}

type CategoryOptions struct {
	CategoryID      uint            `json:"category_id" validate:"required" gorm:"column:category_id"`
	ProductCategory ProductCategory `gorm:"foreignKey:CategoryID"`
	ProductBrand    string          `json:"product_brand" validate:"required" gorm:"varchar(255);not null; column:product_brand"`
}

type Products struct {
	UserID           uint            `gorm:"column:user_id"`
	User             User            `gorm:"foreignKey:UserID"`
	ProductID        uint            `gorm:"primaryKey;autoIncrement"`
	ProductName      string          `gorm:"varchar(255);not null"`
	Categoryid       uint            `gorm:"column:category_id"`
	ProductCategory  ProductCategory `gorm:"foreignKey:Categoryid"`
	ProductBrand     string          `gorm:"varchar(255);not null; column:product_brand"`
	ProductCondition string          `gorm:"varchar(255);not null; conlumn:product_condition"`
	Price            float32         `gorm:"not null; column:price"`
}
