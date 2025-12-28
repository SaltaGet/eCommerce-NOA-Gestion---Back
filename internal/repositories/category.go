package repositories

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"google.golang.org/grpc/metadata"
)

func (repo *CategoryRepository) CategoryGetAll(tenantID string) ([]*pb.Category, error) {
	catReq := &pb.ListCategoriesRequest{}

	ctxt, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	md := metadata.Pairs("x-tenant-identifier", tenantID)

	outCtx := metadata.NewOutgoingContext(ctxt, md)
	resp, err := repo.Client.ListCategories(outCtx, catReq)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return resp.Categories, nil
}