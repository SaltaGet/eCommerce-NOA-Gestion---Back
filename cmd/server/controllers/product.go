package controllers

import (
	"errors"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (ctrl *ProductController) ProductGetByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		log.Error().Msg("El código de consulta está vacío")
		return schemas.ErrorResponse(fiber.StatusBadRequest, "El código es obligatorio", errors.New("código de consulta vacío"))
	}

	tenantID := c.Locals("tenant_identifier").(string)

	product, err := ctrl.ProductService.ProductGetByCode(code, tenantID, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido con éxito",
	})
}

func (ctrl *ProductController) ProductGetPage(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	limit := c.QueryInt("limit")
	if limit == 0 {
		limit = 10
	}
	pageSize := c.QueryInt("page_size")
	if pageSize == 0 {
		pageSize = 10
	}

	req := &schemas.ProductRequest{
		Page:    int32(page),
		PageSize: int32(limit),
		Limit: int32(limit),
	}

	tenantID := c.Locals("tenant_identifier").(string)
	products, total, err := ctrl.ProductService.ProductGetPage(req, tenantID, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(schemas.Response{
		Status:  true,
		Body: fiber.Map{
			"products": products,
			"total":    total,
		},
		Message: "Productos obtenidos con éxito",
	})
}