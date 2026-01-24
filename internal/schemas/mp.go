package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ShoppingCart struct {
	Items  []ItemsCart `json:"items" validate:"required,min=1,dive"`
	Client Payer       `json:"client" validate:"required"`
}

type ItemsCart struct {
	ProductID int64   `json:"product_id" validate:"required" example:"1"`
	Name      string  `json:"name" validate:"required" example:"Product 1"`
	Quantity  int     `json:"quantity" validate:"required,gt=0" example:"1"`
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

type WebhookPayload struct {
	Action string `json:"action"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
}

// additional_info, date_approved, date_created, id, status, status_detail, transaction_amount, transaction_details

type DataInfoPay struct {
	DateApproved       string             `json:"date_approved"`
	DateCreated        string             `json:"date_created"`
	Id                 float64            `json:"id"`
	Status             string             `json:"status"`
	StatusDetail       string             `json:"status_detail"`
	TransactionAmount  float64            `json:"transaction_amount"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
	AditionalInfo      struct {
		Items []ItemsCartPay `json:"items"`
		Payer PayerInfo      `json:"payer"`
	} `json:"additional_info"`
	PayMethod struct {
		Type string `json:"type"`
	} `json:"payment_method"`
	Payer         PayerInfo `json:"payer"`
	OperationType string    `json:"operation_type"`
	Metadata      struct {
		TenantIdentifier string `json:"tenant_identifier"`
		Email            string `json:"email"`
	} `json:"metadata"`
	Message           string `json:"message,omitempty"`
	ExternalReference string `json:"external_reference"`
}

type PayerInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ItemsCartPay struct {
	ProductID string `json:"id"`
	Title     string `json:"title"`
	Quantity  string `json:"quantity"`
	UnitPrice string `json:"unit_price"`
}

type TransactionDetails struct {
	NetReceivedAmount float64 `json:"net_received_amount"`
	TotalPaidAmount   float64 `json:"total_paid_amount"`
}

type PaymentSummary struct {
	// Información básica del pago
	ID                int64     `json:"id"`
	Status            string    `json:"status"`        // ej: "approved"
	StatusDetail      string    `json:"status_detail"` // ej: "accredited"
	DateApproved      time.Time `json:"date_approved"`
	DateCreated			 time.Time `json:"date_created"`
	ExternalReference string    `json:"external_reference"` // Tu ID de orden

	// Datos financieros
	Amount            float64 `json:"transaction_amount"`
	Currency          string  `json:"currency_id"`

	// Método de pago
	PaymentMethodID     string `json:"payment_method_id"`    // ej: "debvisa"
	PaymentType         string `json:"payment_type"`         // ej: "debit_card"
	CardLastFour        string `json:"card_last_four"`       // card.last_four_digits
	StatementDescriptor string `json:"statement_descriptor"` // Lo que aparece en el resumen de la tarjeta

	// El Cliente (Payer)
	PayerEmail          string `json:"payer_email"`
	PayerName           string `json:"payer_name"` // additional_info.payer (first + last name)

	// Proveedor
	TenantIdentifier string `json:"tenant_identifier"`
	TenantEmail      string `json:"tenant_email"`
	TenantPhone			string `json:"tenant_phone"`

	// Detalle del producto (Carrito)
	Items []Item `json:"items"`
}

type Item struct {
	Title     string  `json:"title"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
