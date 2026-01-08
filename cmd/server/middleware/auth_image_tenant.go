package middleware

import (
	"errors"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthImageTenant(c *fiber.Ctx) error {
	token := c.Get("x-token-tenant")
	if token == "" {
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
		return schemas.ErrorResponse(400, "claims no válido", errors.New("claims inválidos"))
	}
	tenantIdentifier, ok := mapClaims["tenant_identifier"].(string)
	if !ok {
		return schemas.ErrorResponse(400, "tenant no válido", errors.New("los tenant no coinciden"))
	}

	tenantID := c.Locals("tenant_identifier")
	if tenantID != tenantIdentifier {
		return schemas.ErrorResponse(400, "tenant no válido", errors.New("los tenant no coinciden"))
	}

	return c.Next()
}