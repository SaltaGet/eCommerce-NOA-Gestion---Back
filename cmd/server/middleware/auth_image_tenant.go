package middleware

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func AuthImageTenant(c *fiber.Ctx) error {
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

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	if tenantID.Identifier != tenantIdentifier {
		log.Error().Msg("tenant no válido")
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
			Status: false,
			Message: "tenant no válido",
		})
	}

	return c.Next()
}