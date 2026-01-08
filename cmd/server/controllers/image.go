package controllers

import (
	"path/filepath"
	"strings"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ImageGet godoc
//
//	@Summary		ImageGet
//	@Description	Obtener imagen
//	@Tags			Image
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Param			filename	path		string	true	"codigo del producto"
//	@Success		200			{object}	any
//	@Router			/ecommerce/{tenantID}/api/v1/image/get/{filename} [get]
func ImageGet(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_identifier").(string)
	filename := c.Params("filename")

	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "el nombre de archivo es obligatorio",
		})
	}

	if filepath.Ext(filename) != ".webp" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "formato de archivo no permitido",
		})
	}

	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	
	var uuidPart string
	if strings.HasSuffix(nameWithoutExt, "p200") {
		uuidPart = strings.TrimSuffix(nameWithoutExt, "p200")
	} else if strings.HasSuffix(nameWithoutExt, "p500") {
		uuidPart = strings.TrimSuffix(nameWithoutExt, "p500")
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "el archivo no tiene un sufijo de tamaño válido (p200/p500)",
		})
	}

	if !utils.IsValidUUIDv7(uuidPart) {
		log.Error().Str("filename", filename).Msg("Intento de acceso con UUID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "el identificador de imagen no es un UUIDv7 válido",
		})
	}

	imagePath, exist := utils.GetPath(tenantID, filename)
	if !exist {
		log.Error().Str("tenant", tenantID).Str("file", filename).Msg("Imagen no encontrada")
		return c.Status(fiber.StatusNotFound).JSON(schemas.Response{
			Status:  false,
			Message: "la imagen no existe",
		})
	}

	return c.SendFile(imagePath)
}
