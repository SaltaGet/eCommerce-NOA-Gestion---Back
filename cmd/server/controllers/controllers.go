package controllers

import "github.com/SaltaGet/ecommerce-fiber-ms/internal/ports"

type TenantController struct {
	TenantService ports.TenantService
}

type ProductController struct {
	ProductService ports.ProductService
}

type CategoryController struct {
	CategoryService ports.CategoryService
}
