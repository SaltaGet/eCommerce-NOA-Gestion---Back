package repositories

import "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"

type ProductRepository struct {
	Client pb.ProductServiceClient	
}

type TenantRepository struct {
	Client pb.TenantServiceClient
}

type CategoryRepository struct {
	Client pb.CategoryServiceClient
}

type MercadoPagoRepository struct {
	Client pb.MPServiceClient
}