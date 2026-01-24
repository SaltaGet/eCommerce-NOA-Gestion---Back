package controllers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// MPInfoPay godoc
//
//	@Summary		MPInfoPay
//	@Description	Obtener información de pago
//	@Tags			Mercado Pago
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Param			reference	path		string	true	"Referencia de pago"
//	@Success		200			{object}	schemas.Response{body=schemas.PaymentSummary}
//	@Router			/ecommerce/{tenantID}/api/v1/mp/info_pay/{reference} [get]
func (m *MPController) MPInfoPay(c *fiber.Ctx) error {
	reference := c.Params("reference")
	if reference == "" {
		log.Error().Msg("Referencia de pago es requerida")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Referencia de pago es requerida",
		})
	}
	ok := utils.IsValidUUIDv4(reference)
	if !ok {
		log.Error().Msg("Referencia de pago no válida")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Referencia de pago no válida",
		})
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	if tenantID.TokenMP == nil {
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Token de Mercado Pago no configurado",
		})
	}

	info, err := m.MPService.MPInfoPay(tenantID, reference)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(info)
}

// MPGenerateLink godoc
//
//	@Summary		MPGenerateLink
//	@Description	Obtener link de pago
//	@Tags			Mercado Pago
//	@Accept			json
//	@Produce		json
//	@Param			tenantID		path		string					true	"ID del Tenant"
//	@Param			shopping_cart	body		schemas.ShoppingCart	true	"Body"
//	@Success		200				{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/mp/generate_link [post]
func (m *MPController) MPGenerateLink(c *fiber.Ctx) error {
	var shoppingCart schemas.ShoppingCart
	if err := c.BodyParser(&shoppingCart); err != nil {
		log.Error().Msg(err.Error())
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}
	if err := shoppingCart.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	if tenantID.TokenMP == nil {
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Token de Mercado Pago no configurado",
		})
	}

	link, err := m.MPService.MPGenerateLink(&shoppingCart, tenantID.Identifier, tenantID.TokenMP, c.BaseURL(), context.Background())
	if err != nil {
		return schemas.HandleError(c, err)
	}
	return c.JSON(link)
}

// MPStatePay godoc
//
//	@Summary		MPStatePay
//	@Description	Obtener link de pago
//	@Tags			Mercado Pago
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/mp/state_pay [post]
func (m *MPController) MPStatePay(c *fiber.Ctx) error {
	body := c.Body()
	var payload schemas.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Error().Msg(err.Error())
		return c.Status(500).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	if tenantID.TokenMP == nil {
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Token de Mercado Pago no configurado",
		})
	}
	// if tenantID.TokenEmail == nil {
	// 	return c.Status(400).JSON(schemas.Response{
	// 		Status:  false,
	// 		Body:    nil,
	// 		Message: "Token de Email no configurado",
	// 	})
	// }
	tenantID.Email = os.Getenv("EMAIL")
	tokenEmail := os.Getenv("TOKEN_EMAIL")
	tenantID.TokenEmail = &tokenEmail

	err := m.MPService.MPStatePay(&payload, &tenantID, tenantID.TokenMP, context.Background(), c)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).Send(nil)
}
