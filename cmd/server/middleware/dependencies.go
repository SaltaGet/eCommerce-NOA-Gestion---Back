package middleware

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func InjectDependencies(deps *dependencies.ContainerGrpc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("dependencies", deps)
		return c.Next()
	}
}
