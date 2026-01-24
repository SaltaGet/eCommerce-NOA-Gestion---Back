package controllers

import (
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ProductGetByCode godoc
//
//	@Summary		ProductGetByCode
//	@Description	Obtener un producto por su código
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Param			code		query		string	true	"codigo del producto"
//	@Success		200			{object}	schemas.Response{body=schemas.ProductResponse}
//	@Router			/ecommerce/{tenantID}/api/v1/product/get_by_code [get]
func (ctrl *ProductController) ProductGetByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		log.Error().Msg("El código de consulta está vacío")
		return schemas.ErrorResponse(fiber.StatusBadRequest, "El código es obligatorio", errors.New("código de consulta vacío"))
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)

	product, err := ctrl.ProductService.ProductGetByCode(code, tenantID.Identifier, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido con éxito",
	})
}

// ProductGetByCode godoc
//
//	@Summary		ProductGetByCode
//	@Description	Obtener un producto por su código
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Param			page		query		int		false	"pagina"
//	@Param			limit		query		int		false	"limite"
//	@Param			name	query		string		false	"nombre del producto"
//	@Param			category_id	query		int		false	"categoria id del producto"
//	@Param			order	query		string		false	"ordenar por" Enums(PRICE_LOW_TO_HIGH, PRICE_HIGH_TO_LOW, NAME_A_Z, NAME_Z_A)
//	@Success		200			{object}	schemas.Response{body=schemas.ProductResponseDTO}
//	@Router			/ecommerce/{tenantID}/api/v1/product/get_page [get]
func (ctrl *ProductController) ProductGetPage(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	limit := c.QueryInt("limit")
	if limit == 0 {
		limit = 10
	}
	category, err := strconv.ParseInt(c.Query("category_id", "0"), 10, 32)
	name := c.Query("name")

	var search *string
	if name != "" {
		search = &name
	}
	var categoryID *int32
	if category != 0 {
		val := int32(category)
    categoryID = &val
	}

	order := c.Query("order", "")
	var orderBy *schemas.SortBy 
	if order != "" {
		if val, ok := schemas.SortBy_value[order]; ok {
        sortVal := schemas.SortBy(val)
        orderBy = &sortVal
    }
	}

	req := &schemas.ProductRequest{
		Page:     int32(page),
		Limit:    int32(limit),
		CategoryID: categoryID,
		Search:     search,
		Sort: orderBy,
	}

	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	products, total, err := ctrl.ProductService.ProductGetPage(req, tenantID.Identifier, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": products, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Productos obtenidos correctamente",
	})
}

// ProductSaveImage godoc
//
//	@Summary		ProductSaveImage
//
//	@Description	### Flujo de Carga de Imágenes
//	@Description	Genera un token temporal desde la API principal para subir imágenes al microservicio.
//	@Description
//	@Description	**Endpoitn de la api princial para pedir token:**
//	@Description	~~~
//	@Description	POST /api/v1/product/generate_token_to_image
//	@Description	~~~
//	@Description
//	@Description	**Pasos requeridos:**
//	@Description	1. Incluir el token en el header `x-token-tenant`.
//	@Description
//	@Description	> *Nota: El token tiene una validez limitada de 30 minutos.*
//
//	@Tags			Product
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			tenantID		path		string	true	"ID del Tenant"
//	@Param			x-token-tenant	header		string	true	"Token o ID de validación personalizada"
//	@Param			primaryImage	formData	file	false	"Imagen principal del producto"
//	@Param			secondaryImage	formData	[]file	false	"Imágenes secundarias (puedes enviar varias)"
//	@Success		200				{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/product/upload_image [post]
func (ctrl *ProductController) ProductSaveImage(c *fiber.Ctx) error {
	primaryImage, err := c.FormFile("primaryImage")
	if err != nil {
		primaryImage = nil
	}

	var imageFiles []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if files, ok := form.File["secondaryImage"]; ok {
			imageFiles = files
		}
	}

	schema := &schemas.ProductUploadSchema{
		PrimaryImage:    primaryImage,
		SecondaryImages: imageFiles,
	}

	add := c.Locals("add").(float64)
	lenImgs := len(imageFiles)
	// if lenImgs == int(add) {
	// 	return c.Status(fiber.StatusOK).JSON(schemas.Response{
	// 		Status:  false,
	// 		Message: "Las imágenes secundarias son obligatorias",
	// 	})
	// }
	tenantID := c.Locals("tenant_data").(schemas.TenantResponseSetting)
	productID := c.Locals("product_id").(float64)
	keep := c.Locals("keep").(string)
	remove := c.Locals("remove").(string)
	primaryImg := c.Locals("primary_image").(string)
	if primaryImg == "set" && primaryImage == nil {
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Message: "La imagen principal no debede de estar vacio, no concuerda con el token de validación",
		})
	}

	
	if lenImgs != int(add) {
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Message: "Las imágenes secundarias para agregar no concuerdan con la cantidad de imagenes a agregar",
		})
	}

	addInt := int32(add)

	validationData := &schemas.ProductValidateImage{
		ProductID:    int64(productID),
		PrimaryImage: primaryImg,
		SecondaryImage: schemas.ValidateSecondaryImage{
			Add:         &addInt,
			KeepUUIDs:   utils.SplitStrings(&keep),
			RemoveUUIDs: utils.SplitStrings(&remove),
		},
	}
	if err := validationData.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err = ctrl.ProductService.ProductUploadImages(tenantID.Identifier, schema, int64(productID), validationData, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Imágenes de producto guardadas con éxito",
	})
}

// func (ctrl *ProductController) ProductSaveImage(c *fiber.Ctx) error {
// 	token	:= c.FormValue("token")
// 	if token == "" {
// 		log.Error().Msg("El token de validación está vacío")
// 	}

// 	primaryImage, err := c.FormFile("primaryImage")
// 	if err != nil {
// 		log.Error().Msg("La imagen principal es obligatoria")
// 		return schemas.ErrorResponse(fiber.StatusBadRequest, "La imagen principal es obligatoria", errors.New("imagen principal vacía"))
// 	}

// 	secondaryImages, err := c.MultipartForm()
// 	if err != nil {
// 		log.Error().Msg("Error al obtener las imágenes secundarias")
// 		return schemas.ErrorResponse(fiber.StatusBadRequest, "Error al obtener las imágenes secundarias", err)
// 	}

// 	imageFiles := make([]*multipart.FileHeader, 0)
// 	if secondaryImages != nil {
// 		if files, ok := secondaryImages.File["secondaryImage"]; ok {
// 			for _, file := range files {
// 				imageFiles = append(imageFiles, file)
// 			}
// 		}
// 	}

// 	schema := &schemas.ProductUploadSchema{
// 		Token:          token,
// 		PrimaryImage:   primaryImage,
// 		SecondaryImages: imageFiles,
// 	}

// 	tenantID := c.Locals("tenant_identifier").(string)

// 	err = ctrl.ProductService.ProductUploadImages(tenantID, schema, c.Context())
// 	if err != nil {
// 		return schemas.HandleError(c, err)
// 	}

// 	return c.JSON(schemas.Response{
// 		Status:  true,
// 		Message: "Imágenes de producto guardadas con éxito",
// 	})
// }
