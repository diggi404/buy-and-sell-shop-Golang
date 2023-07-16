package models

import (
	"time"
)

type User struct {
	ID                   uint          `json:"-" gorm:"primaryKey;autoIncrement"`
	Name                 string        `gorm:"type:varchar(100);not null"`
	Email                string        `gorm:"type:varchar(255);not null;unique"`
	Password             string        `json:"-" gorm:"type:varchar(255);not null"`
	Closed               bool          `gorm:"not null;default:false"`
	DefaultPaymentMethod uint          `json:"default_payment_method" gorm:"column:default_payment_method;unique"`
	ResetRequested       bool          `json:"reset_requested" gorm:"column:reset_requested;default:false"`
	EmailVerified        bool          `json:"email_verified" gorm:"column:email_verified;not null;default:false"`
	CreatedAt            time.Time     `json:"-" gorm:"timestamp;not null"`
	UpdatedAt            time.Time     `json:"-" gorm:"timestamp;not null"`
	CreditCards          []CreditCard  `json:"credit_cards"`
	Momo                 []MobileMoney `json:"momo"`
}

type EmailVerify struct {
	UserId    uint      `json:"user_id" gorm:"column:user_id;primaryKey"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
	Link      string    `json:"link" gorm:"column:link;not null;unique"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	ExpiresAt time.Time `json:"-" gorm:"timestamp;not null"`
}

type AddressBook struct {
	AddressId uint      `gorm:"primaryKey;autoIncrement; column:address_id"`
	UserId    uint      `json:"-" gorm:"column:user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
	FirstName string    `gorm:"varchar(100);not null; column:fname"`
	LastName  string    `gorm:"varchar(100);not null; column:lname"`
	Address1  string    `gorm:"varchar(255);not null"`
	City      string    `gorm:"varchar(100);not null"`
	State     string    `gorm:"varchar(100);notn null"`
	ZipCode   string    `gorm:"varchar(50);not null"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt time.Time `json:"-" gorm:"timestamp;not null"`
}

type ProductCategory struct {
	CategoryId   uint       `gorm:"primaryKey;autoIncrement; column:category_id"`
	CategoryName string     `json:"category_name" gorm:"varchar(255);not null; column:category_name"`
	CreatedAt    time.Time  `json:"-" gorm:"timestamp;not null"`
	UpdatedAt    time.Time  `json:"-" gorm:"timestamp;not null"`
	Products     []Products `json:"products" gorm:"foreignKey:Categoryid"`
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
	UserID           uint            `json:"seller_id" gorm:"column:user_id"`
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
	ItemLikes        int             `json:"likes" gorm:"column:likes;not null;default:0"`
	CreatedAt        time.Time       `json:"-" gorm:"timestamp;not null"`
	UpdatedAt        time.Time       `json:"-" gorm:"timestamp;not null"`
}

type Cart struct {
	CartItemId uint      `json:"-" gorm:"primaryKey;autoIncrement;column:cart_item_id"`
	Userid     uint      `json:"-" gorm:"column:user_id"`
	User       User      `json:"-" gorm:"foreignKey:Userid"`
	ProductId  uint      `json:"product_id" gorm:"column:product_id"`
	Products   Products  `json:"product" gorm:"foreignKey:ProductId;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt  time.Time `json:"-" gorm:"timestamp;not null"`
}

type TotalCart struct {
	CartID     uint      `json:"cart_id" gorm:"primaryKey;column:cart_id"`
	User_Id    uint      `json:"user_id" gorm:"column:user_id;not null;unique"`
	NewUser    User      `json:"-" gorm:"foreignKey:User_Id"`
	TotalPrice float32   `json:"total_price" gorm:"column:total_price"`
	CreatedAt  time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt  time.Time `json:"-" gorm:"timestamp;not null"`
}

type MobileMoney struct {
	NumberId  uint      `json:"number_id" gorm:"primaryKey;autoIncrement;column:number_id"`
	User_ID   uint      `json:"-" gorm:"column:user_id"`
	User      User      `json:"-" gorm:"foreignKey:User_ID"`
	Number    uint      `json:"number" gorm:"column:number;not null" validate:"required,min=10,max=10"`
	Network   string    `json:"network" gorm:"column:network;not null" validate:"required"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt time.Time `json:"-" gorm:"timestamp;not null"`
}

type BillingAddress struct {
	AddressiD uint      `json:"-" gorm:"column:address_id;primaryKey;autoIncrement"`
	UserId    uint      `json:"-" gorm:"column:user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
	FirstName string    `json:"fname" gorm:"column:fname;not null"`
	LastName  string    `json:"lname" gorm:"column:lname;not null"`
	Address1  string    `json:"address1" gorm:"column:address1;not null"`
	City      string    `json:"city" gorm:"column:city;not null"`
	State     string    `json:"state" gorm:"column:state;notn null"`
	ZipCode   string    `json:"zip_code" gorm:"column:zip_code;not null"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt time.Time `json:"-" gorm:"timestamp;not null"`
}

type CreditCard struct {
	CardId       uint           `json:"card_id" gorm:"primaryKey;autoIncrement;column:card_id"`
	User_ID      uint           `json:"-" gorm:"column:user_id;not null"`
	User         User           `json:"-" gorm:"foreignKey:User_ID"`
	AddressID    uint           `json:"address_id" gorm:"column:address_id;not null"`
	CardNumber   uint           `json:"-" gorm:"column:card_number;not null"`
	MaskedNumber string         `json:"masked_number" gorm:"column:masked_number;not null"`
	CardMonth    uint           `json:"card_month" gorm:"column:card_month;not null"`
	CardYear     uint           `json:"card_year" gorm:"column:card_year;not null"`
	CardType     string         `json:"card_type" gorm:"column:card_type; not null"`
	LastFour     uint           `json:"last_four" gorm:"column:last_four;not null"`
	Address      BillingAddress `json:"billing_address" gorm:"foreignKey:AddressID"`
	IsDefault    bool           `json:"is_default" gorm:"column:is_default;default:false"`
	CreatedAt    time.Time      `json:"-" gorm:"timestamp;not null"`
	UpdatedAt    time.Time      `json:"-" gorm:"timestamp;not null"`
}

type Shipment struct {
	TrackingId     uint           `json:"tracking_id" gorm:"primaryKey;autoIncrement;column:tracking_id"`
	ItemId         uint           `json:"item_id" gorm:"column:item_id;unique"`
	PurchasedItems PurchasedItems `json:"-" gorm:"foreignKey:ItemId"`
	TrackingNumber string         `json:"tracking_number" gorm:"column:tracking_number" validate:"required"`
	Carrier        string         `json:"carrier" gorm:"column:carrier" validate:"required"`
}

type PurchasedItems struct {
	ItemID           uint       `json:"item_id" gorm:"primaryKey;autoIncrement;column:item_id"`
	OrderID          uint       `json:"order_id" gorm:"column:order_id"`
	Orders           Orders     `json:"-" gorm:"foreignKey:OrderID"`
	SellerId         uint       `json:"seller_id" gorm:"column:seller_id"`
	ProductName      string     `json:"product_name" gorm:"not null"`
	ProductBrand     string     `json:"product_brand" gorm:"not null; column:product_brand"`
	ProductCondition string     `json:"product_condition" gorm:"not null; conlumn:product_condition"`
	ShoeSize         float32    `json:"shoe_size,omitempty" gorm:"column:shoe_size"`
	ClothSize        string     `json:"cloth_size,omitempty" gorm:"column:cloth_size"`
	Color            string     `json:"color,omitempty" gorm:"column:color"`
	Price            float32    `json:"price" gorm:"not null; column:price"`
	OrderStatus      string     `json:"order_status" gorm:"column:order_status;default:processing"`
	CanCancel        bool       `json:"can_cancel" gorm:"column:can_cancel;default:true"`
	Seller           Sellers    `json:"seller_info" gorm:"foreignKey:SellerId"`
	Shipment         []Shipment `json:"shipment" gorm:"foreignKey:ItemId"`
	CreatedAt        time.Time  `json:"-" gorm:"timestamp;not null"`
	UpdatedAt        time.Time  `json:"-" gorm:"timestamp;not null"`
}

type Orders struct {
	OrderId        uint             `json:"order_id" gorm:"primaryKey;autoIncrement;column:order_id"`
	UserId         uint             `json:"-" gorm:"column:user_id"`
	User           User             `json:"-" gorm:"foreignKey:UserId"`
	PaymentId      uint             `json:"payment_id" gorm:"column:payment_id"`
	PaymentMethod  string           `json:"payment_method" gorm:"column:payment_method"`
	PaidTotal      float32          `json:"paid_total" gorm:"column:paid_total;not null"`
	PurchasedItems []PurchasedItems `json:"items" gorm:"foreignKey:OrderID"`
	AddressID      uint             `json:"-" gorm:"column:address_id"`
	AddressBook    AddressBook      `json:"shipping_address" gorm:"foreignKey:AddressID"`
	CreatedAt      time.Time        `json:"-" gorm:"timestamp;not null"`
	UpdatedAt      time.Time        `json:"-" gorm:"timestamp;not null"`
}

type Sellers struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey;column:user_id"`
	Seller    User      `json:"-" gorm:"foreignKey:UserID"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email;unique"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	UpdatedAt time.Time `json:"-" gorm:"timestamp;not null"`
}

type OtpEmail struct {
	UserId    uint      `json:"user_id" gorm:"primaryKey;column:user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
	Value     int       `json:"value" gorm:"column:value;not null"`
	Email     string    `json:"email" gorm:"column:email;not null"`
	CreatedAt time.Time `json:"-" gorm:"timestamp;not null"`
	ExpiresAt time.Time `json:"-" gorm:"timestamp;not null"`
}
