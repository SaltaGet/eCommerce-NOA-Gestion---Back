package ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

type TenantRepository interface {
	TenantList() (*pb.ListTenantsResponse, error)
	TenantGet(tenantIdentifier string) (*pb.TenantResponse,error)
	TenantSaveImage(req *pb.TenantRequestImageSetting, ctx context.Context) (*pb.TenantUpdateImageResponse, error)
}

type TenantService interface {
	TenantList() ([]schemas.TenantResponse, error)
	TenantGet(tenantIdentifier string) (*schemas.TenantResponseSetting, error)
	TenantSaveImage(tenantID string, schema *schemas.TenantUploadSchema, ctx context.Context) error
}
