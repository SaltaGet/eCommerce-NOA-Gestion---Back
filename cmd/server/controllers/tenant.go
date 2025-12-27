package controllers

import "github.com/gofiber/fiber/v2"

func (c *TenantController) TenantGetAll(ctx *fiber.Ctx) error {
	tenants, err := c.TenantService.TenantList()
	if err != nil {
		return err
	}
	return ctx.JSON(tenants)
}