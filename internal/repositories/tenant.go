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

func (r *TenantRepository) TenantGet(identifier string) (*pb.TenantResponse, error) {
	req := &pb.TenantRequest{
		Identifier: identifier,
	}
	tenant, err := r.Client.TenantGetIdentifier(context.Background(), req)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return tenant, nil
}

func (r *TenantRepository) TenantSaveImage(req *pb.TenantRequestImageSetting, ctx context.Context) (*pb.TenantUpdateImageResponse, error) {
	reps, err := r.Client.TenantUpdateImageSetting(ctx, req)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return reps, nil
}