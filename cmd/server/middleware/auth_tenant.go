package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthTenantMiddleware(c *fiber.Ctx) error {
	host := c.Hostname()

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Host inválido para determinar tenant",
		})
	}

	for _, tenantID := range []string{"tenant1", "pepe", "tenant3"} {
		if parts[0] == tenantID {
			return c.Next()
		}
	}

	// tenantIdentifier := parts[0]
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Tenant no autorizado",
	})
}