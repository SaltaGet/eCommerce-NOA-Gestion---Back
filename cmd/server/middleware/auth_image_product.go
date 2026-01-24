package middleware

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func AuthImageProduct(c *fiber.Ctx) error {
	token := c.Get("x-token-tenant")
	if token == "" {
		log.Error().Msg("Token de validación no proporcionado")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "Token de validación no proporcionado",
		})
	}

	claims, err := utils.VerifyToken(token)
	if err != nil {
		return schemas.HandleError(c, err)
	}
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		log.Error().Msg("claims no válidos")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "claims no válidos",
		})
	}
	tenantIdentifier, ok := mapClaims["tenant_identifier"].(string)
	if !ok {
		log.Error().Msg("tenant no válido, claim identifier")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim identifier",
		})
	}
	productID, ok := mapClaims["product_id"].(float64)
	if !ok {
		log.Error().Msg("tenant no válido, claim product_id")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim product_id",
		})
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	if tenantID.Identifier != tenantIdentifier {
		log.Error().Msg("tenant no válido")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido",
		})
	}

	keep, ok := mapClaims["keep"].(string)
	if !ok {
		log.Error().Msg("tenant no válido, claim keep")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim keep",
		})
	}
	remove, ok := mapClaims["remove"].(string)
	if !ok {
		log.Error().Msg("tenant no válido, claim remove")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim remove",
		})
	}
	primaryImage, ok := mapClaims["primary_image"].(string)
	if !ok {
		log.Error().Msg("tenant no válido, claim primary")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim primary",
		})
	}
	add, ok := mapClaims["add"].(float64)
	if !ok {
		log.Error().Msg("tenant no válido, claim add")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido, claim add",
		})
	}
	
	c.Locals("product_id", productID)

	c.Locals("keep", keep)
	c.Locals("remove", remove)
	c.Locals("primary_image", primaryImage)
	c.Locals("add", add)

	return c.Next()
}