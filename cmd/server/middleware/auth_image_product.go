package middleware

import (
	"errors"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthImageProduct(c *fiber.Ctx) error {
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
	productID, ok := mapClaims["product_id"].(float64)
	if !ok {
		return schemas.ErrorResponse(400, "product_id no válido", errors.New("product_id inválido"))
	}

	tenantID := c.Locals("tenant_identifier")
	if tenantID != tenantIdentifier {
		return schemas.ErrorResponse(400, "tenant no válido", errors.New("los tenant no coinciden"))
	}

	keep, ok := mapClaims["keep"].(string)
	if !ok {
		return schemas.ErrorResponse(400, "claims no válido", errors.New("keep inválido"))
	}
	remove, ok := mapClaims["remove"].(string)
	if !ok {
		return schemas.ErrorResponse(400, "claims no válido", errors.New("remove inválido"))
	}
	primaryImage, ok := mapClaims["primary_image"].(string)
	if !ok {
		return schemas.ErrorResponse(400, "claims no válido", errors.New("imagen principal no valida"))
	}
	add, ok := mapClaims["add"].(float64)
	if !ok {
		return schemas.ErrorResponse(400, "claims no válido", errors.New("add no valido"))
	}
	
	c.Locals("product_id", productID)

	c.Locals("keep", keep)
	c.Locals("remove", remove)
	c.Locals("primary_image", primaryImage)
	c.Locals("add", add)

	return c.Next()
}