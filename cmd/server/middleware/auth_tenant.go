package middleware

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
  "github.com/rs/zerolog/log"
)

func AuthTenantMiddleware(store *cache.TenantStore) fiber.Handler {
  return func(c *fiber.Ctx) error {
    tenantID := c.Params("tenantID")
    
    tenant, exists := store.Get(tenantID)
    if !exists {
      log.Warn().Str("tenant_id", tenantID).Msg("Intento de acceso con tenant no autorizado")
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Tenant no autorizado, verifique nuevamente. Si cree que se trata de un error contacte al administrador",
      })
    }

    c.Locals("tenant_data", tenant) // Ahora tienes todo el objeto, no solo el ID
    return c.Next()
  }
}

// func AuthTenantMiddleware(c *fiber.Ctx) error {
// 	tenantID := c.Params("tenantID")
// 	if tenantID == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Tenant ID no proporcionado",
// 		})
// 	}

	

// 	for _, tenant := range tenants {
// 		if tenantID == tenant.Identifier {
// 			c.Locals("tenant_identifier", tenantID)
// 			return c.Next()
// 		}
// 	}

// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 		"error": "Tenant no autorizado, verifique nuevamente. Si cree que se trata de un error contacte al administrador",
// 	})
// }
