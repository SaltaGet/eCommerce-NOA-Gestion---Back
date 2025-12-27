package repositories

import (
	"context"
	"time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"google.golang.org/grpc/metadata"
)

const requestTimeout = 5 * time.Second

func (repo *ProductRepository) ProductGetByCode(code string, tenantID string, ctx context.Context) (*pb.Product, error) {
	prodReq := &pb.GetProductRequest{
		Code: code,
	}

	ctxt, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	md := metadata.Pairs("x-tenant-identifier", tenantID)

	outCtx := metadata.NewOutgoingContext(ctxt, md)
	return repo.Client.GetProduct(outCtx, prodReq)
}

func (repo *ProductRepository) ProductGetPage(req *schemas.ProductRequest, tenantID string, ctx context.Context) (*pb.ListProductsResponse, error) {
	prodReq := &pb.ListProductsRequest{
		Page:     req.Page,
		Limit:    req.Limit,
		PageSize: req.Limit,
	}

	if req.Search != nil {
		prodReq.Search = req.Search
	}

	if req.CategoryID != nil {
		var cat int32 = 0
		if req.CategoryID == &cat {
			prodReq.CategoryId = nil
		} else {
			prodReq.CategoryId = req.CategoryID
		}
	}

	if req.Search != nil {
		prodReq.Search = req.Search
	}

	if req.Sort != nil {
		prodReq.Sort = (*pb.ListProductsRequest_SortBy)(req.Sort)
	}

	ctxt, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	md := metadata.Pairs("x-tenant-identifier", tenantID)

	outCtx := metadata.NewOutgoingContext(ctxt, md)
	return repo.Client.ListProducts(outCtx, prodReq)
}
