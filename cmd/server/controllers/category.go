package controllers

import "github.com/gofiber/fiber/v2"

func (c *CategoryController) CategoryGetAll(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_identifier").(string)

	categories, err := c.CategoryService.CategoryGetAll(tenantID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve categories",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":     false,
		"categories": categories,
	})
}