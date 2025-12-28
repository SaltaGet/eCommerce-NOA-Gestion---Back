package repositories

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

func (r *TenantRepository) TenantList() (*pb.ListTenantsResponse, error) {
	req := &pb.ListTenantsRequest{}
	tenants, err := r.Client.ListTenants(context.Background(), req)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return tenants, nil
}