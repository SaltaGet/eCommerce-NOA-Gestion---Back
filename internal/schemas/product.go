package schemas

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID             int64            `json:"id"`
	Code           string           `json:"code"`
	Name           string           `json:"name"`
	Description    *string          `json:"description"`
	Category       CategoryResponse `json:"category"`
	Price          float64          `json:"price"`
	Stock          float64          `json:"stock"`
	PrimaryImage   *string          `json:"primary_image"`
	SecondaryImage []string         `json:"secondary_images"`
}

type ProductResponseDTO struct {
	ID           int64            `json:"id"`
	Code         string           `json:"code"`
	Name         string           `json:"name"`
	Category     CategoryResponse `json:"category"`
	Price        float64          `json:"price"`
	Stock        float64          `json:"stock"`
	PrimaryImage *string          `json:"primary_image"`
}

type SortBy int32

const (
	ListProductsRequest_PRICE_LOW_TO_HIGH SortBy = 0
	ListProductsRequest_PRICE_HIGH_TO_LOW SortBy = 1
	ListProductsRequest_NAME_A_Z          SortBy = 2
	ListProductsRequest_NAME_Z_A          SortBy = 3
)

// Mapa para obtener el nombre en string (generado automáticamente)
var SortBy_name = map[int32]string{
	0: "PRICE_LOW_TO_HIGH",
	1: "PRICE_HIGH_TO_LOW",
	2: "NAME_A_Z",
	3: "NAME_Z_A",
}

var SortBy_value = map[string]int32{
	"PRICE_LOW_TO_HIGH": 0,
	"PRICE_HIGH_TO_LOW": 1,
	"NAME_A_Z":          2,
	"NAME_Z_A":          3,
}

type ProductRequest struct {
	// Agregamos 'query' y validaciones básicas
	Page     int32 `json:"page" query:"page" validate:"min=1"`
	PageSize int32 `json:"page_size" query:"page_size" validate:"min=1,max=100"`
	Limit    int32 `json:"limit" query:"limit" validate:"min=1,max=100"`
	// Punteros para manejar nulls
	CategoryID *int32  `json:"category_id" query:"category_id"`
	Search     *string `json:"search" query:"search"`
	Sort       *SortBy `json:"sort" query:"sort"`
}

type ProductValidateImage struct {
	ProductID      int64                  `json:"product_id" validate:"required" example:"1"`
	PrimaryImage   string                 `json:"primary_image" validate:"required,oneof=set keep" example:"set | keep"`
	SecondaryImage ValidateSecondaryImage `json:"secondary_image" validate:"required"`
}

type ValidateSecondaryImage struct {
	Add         *int32    `json:"add" example:"1"`
	KeepUUIDs   []string `json:"keep_uuid" example:"lista de uuids que se desea retener"`
	RemoveUUIDs []string `json:"remove_uuid" example:"lista de uuids que se desea remover"`
}

func (p *ProductValidateImage) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type ProductUploadSchema struct {
	PrimaryImage    *multipart.FileHeader   `form:"primaryImage"`
	SecondaryImages []*multipart.FileHeader `form:"secondaryImage"`
}
