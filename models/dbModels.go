package models

import "time"

type User struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(255);not null;unique"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	CartId    uint      `gorm:"column:cart_id"`
	Closed    bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"timestamp;not null"`
	UpdatedAt time.Time `gorm:"timestamp;not null"`
}

type AddressBook struct {
	AddressId uint      `gorm:"primaryKey;autoIncrement; column:address_id"`
	UserId    uint      `gorm:"column:user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
	FirstName string    `gorm:"varchar(100);not null; column:fname"`
	LastName  string    `gorm:"varchar(100);not null; column:lname"`
	Address1  string    `gorm:"varchar(255);not null"`
	City      string    `gorm:"varchar(100);not null"`
	State     string    `gorm:"varchar(100);notn null"`
	ZipCode   string    `gorm:"varchar(50);not null"`
	CreatedAt time.Time `gorm:"timestamp;not null"`
	UpdatedAt time.Time `gorm:"timestamp;not null"`
}

type ProductCategory struct {
	CategoryId   uint      `gorm:"primaryKey;autoIncrement; column:category_id"`
	CategoryName string    `json:"category_name" gorm:"varchar(255);not null; column:category_name"`
	CreatedAt    time.Time `gorm:"timestamp;not null"`
	UpdatedAt    time.Time `gorm:"timestamp;not null"`
}

type CategoryOptions struct {
	CategoryID      uint            `json:"category_id" validate:"required" gorm:"column:category_id"`
	ProductCategory ProductCategory `json:"-" gorm:"foreignKey:CategoryID"`
	ProductBrand    string          `json:"product_brand" validate:"required" gorm:"varchar(255);not null; column:product_brand"`
	ShoeSize        float32         `json:"shoe_size,omitempty" gorm:"column:shoe_size"`
	ClothSize       string          `json:"cloth_size,omitempty" gorm:"varchar(50); column:cloth_size"`
	Color           string          `json:"color,omitempty" gorm:"varchar(100); column:color"`
	CreatedAt       time.Time       `gorm:"timestamp;not null"`
	UpdatedAt       time.Time       `gorm:"timestamp;not null"`
}

type Products struct {
	UserID           uint            `json:"-" gorm:"column:user_id"`
	User             User            `json:"-" gorm:"foreignKey:UserID"`
	ProductID        uint            `json:"product_id" gorm:"primaryKey;autoIncrement;column:product_id"`
	ProductName      string          `json:"product_name" gorm:"varchar(255);not null"`
	Categoryid       uint            `json:"category_id" gorm:"column:category_id"`
	ProductCategory  ProductCategory `json:"-" gorm:"foreignKey:Categoryid"`
	ProductBrand     string          `json:"product_brand" gorm:"varchar(255);not null; column:product_brand"`
	ProductCondition string          `json:"product_condition" gorm:"varchar(255);not null; conlumn:product_condition"`
	ShoeSize         float32         `json:"shoe_size,omitempty" gorm:"column:shoe_size"`
	ClothSize        string          `json:"cloth_size,omitempty" gorm:"varchar(50); column:cloth_size"`
	Color            string          `json:"color,omitempty" gorm:"varchar(100); column:color"`
	Price            float32         `json:"price" gorm:"not null; column:price"`
	CreatedAt        time.Time       `json:"-" gorm:"timestamp;not null"`
	UpdatedAt        time.Time       `json:"-" gorm:"timestamp;not null"`
}

type Cart struct {
	CartItemId       uint      `json:"-" gorm:"primaryKey;autoIncrement;column:cart_id"`
	Userid           uint      `json:"-" gorm:"column:user_id"`
	User             User      `json:"-" gorm:"foreignKey:Userid"`
	ProductId        uint      `json:"product_id" gorm:"column:product_id"`
	Products         Products  `json:"-" gorm:"foreignKey:ProductId"`
	ProductName      string    `json:"product_name" gorm:"varchar(255);not null"`
	ProductBrand     string    `json:"product_brand" gorm:"varchar(255);not null; column:product_brand"`
	ProductCondition string    `json:"product_condition" gorm:"varchar(255);not null; conlumn:product_condition"`
	ShoeSize         float32   `json:"shoe_size,omitempty" gorm:"column:shoe_size"`
	ClothSize        string    `json:"cloth_size,omitempty" gorm:"varchar(50); column:cloth_size"`
	Color            string    `json:"color,omitempty" gorm:"varchar(100); column:color"`
	Price            float32   `json:"price" gorm:"not null; column:price"`
	CreatedAt        time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt        time.Time `json:"-" gorm:"timestamp;not null"`
}

type MobileMoney struct {
	NumberId uint   `json:"number_id" gorm:"primaryKey;autoIncrement;column:number_id"`
	User_ID  uint   `json:"-" gorm:"column:user_id"`
	User     User   `json:"-" gorm:"foreignKey:User_ID"`
	Number   uint   `json:"number" gorm:"column:number" validate:"required,min=10,max=10"`
	Network  string `json:"network" gorm:"column:network" validate:"required"`
}

type CreditCard struct {
	CardId     uint `json:"card_id" gorm:"primaryKey;autoIncrement;column:card_id"`
	User_id    uint `json:"-" gorm:"column:user_id"`
	User       User `json:"-" gorm:"foreignKey:User_id"`
	CardNumber uint `json:"card_number" gorm:"column:card_number" validate:"required,credit_card"`
	CardMonth  uint `json:"card_month" gorm:"column:card_month" validate:"required,min=2,max=2"`
	CardYear   uint `json:"card_year" gorm:"column:card_year" validate:"required,min=4,max=4"`
}
