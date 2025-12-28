package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthTenantMiddleware(c *fiber.Ctx) error {
	tenantID := c.Params("tenantID")
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID no proporcionado",
		})
	}

	for _, tenant := range []string{"tenant5", "daniel", "test"} {
		if tenantID == tenant {
			c.Locals("tenant_identifier", tenantID)
			return c.Next()
		}
	}
	// host := c.Hostname()

	// parts := strings.Split(host, ".")
	// if len(parts) < 2 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Host inválido para determinar tenant",
	// 	})
	// }

	// for _, tenantID := range []string{"tenant5", "daniel", "test"} {
	// 	if parts[0] == tenantID {
	// 		c.Locals("tenant_identifier", tenantID)
	// 		return c.Next()
	// 	}
	// }

	// tenantIdentifier := parts[0]
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Tenant no autorizado",
	})
}