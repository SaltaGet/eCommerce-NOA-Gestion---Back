package services

import "github.com/SaltaGet/ecommerce-fiber-ms/internal/ports"

type TenantService struct {
	Repo ports.TenantRepository
}

type ProductService struct {
	Repo ports.ProductRepository
}

type CategoryService struct {
	Repo ports.CategoryRepository
}

type MPService struct {
	Repo ports.MPRepository
}