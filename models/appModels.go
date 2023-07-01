package models

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=45"`
}

type SignupSchema struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email,lowercase"`
	Password        string `json:"password" validate:"required,min=8,max=45"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type CreateAddressBook struct {
	FirstName string `json:"fname" validate:"required"`
	LastName  string `json:"lname" validate:"required"`
	Address1  string `json:"address1" validate:"required"`
	Address2  string `json:"adress2"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required,min=2,max=2"`
	ZipCode   string `json:"zip_code" validate:"required,min=5,max=5"`
}

type UpdateAddressBook struct {
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Address1  string `json:"address1"`
	City      string `json:"city"`
	State     string `json:"state" validate:"min=2,max=2"`
	ZipCode   string `json:"zip_code" validate:"min=5,max=5"`
}

// type AddProductCategory struct {
// 	Name string `json:"name"`
// }

type AddProductCategory struct {
	Categories []ProductCategory `json:"categories"`
}

type AddCategoryOptions struct {
	Options []CategoryOptions `json:"options" validate:"dive"`
}

type AddProduct struct {
	UserId           uint    `json:"user_id"`
	ProductName      string  `json:"product_name" validate:"required"`
	Categoryid       uint    `json:"category_id" validate:"required"`
	ProductBrand     string  `json:"product_brand" validate:"required"`
	ProductCondition string  `json:"product_condition" validate:"required"`
	Price            float32 `json:"price" validate:"required"`
}
