package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ShoppingCart struct {
	Items  []ItemsCart `json:"items" validate:"required,min=1,dive"`
	Client Payer       `json:"client" validate:"required"`
}

type ItemsCart struct {
	ProductID int64   `json:"product_id" validate:"required" example:"1"`
	Name      string  `json:"name" validate:"required" example:"Product 1"`
	Quantity  int `json:"quantity" validate:"required,gt=0" example:"1"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0" example:"100.00"`
}

type Payer struct {
	Email   string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Name    string `json:"name" validate:"required" example:"John"`
	Surname string `json:"surname" validate:"required" example:"Doe"`
}

func (m *ShoppingCart) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
	if err == nil {
		return nil
	}

	validatorErr := err.(validator.ValidationErrors)[0]
	field := validatorErr.Field()
	tag := validatorErr.Tag()
	params := validatorErr.Param()

	errorMessage := field + " " + tag + " " + params
	return ErrorResponse(422, fmt.Sprintf("error al validar campo(s): %s", errorMessage), err)
}
