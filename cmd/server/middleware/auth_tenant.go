package middleware

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func AuthTenantMiddleware(c *fiber.Ctx) error {
	tenantID := c.Params("tenantID")
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID no proporcionado",
		})
	}

	deps := c.Locals("dependencies").(*dependencies.ContainerGrpc)
	if deps == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Dependencias no proporcionadas",
		})
	}

	tenants, err := deps.Services.TenantService.TenantList()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener tenants",
		})
	}

	for _, tenant := range tenants {
		if tenantID == tenant.Identifier {
			c.Locals("tenant_identifier", tenantID)
			return c.Next()
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Tenant no autorizado",
	})
}
