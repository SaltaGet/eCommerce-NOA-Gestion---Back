package repositories

import (
	"context"
	"mime/multipart"
	"sync"
	"time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/rs/zerolog/log"
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
	prod, err := repo.Client.GetProduct(outCtx, prodReq)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return prod, nil
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
	listProd, err := repo.Client.ListProducts(outCtx, prodReq)
	if err != nil {
		return nil, schemas.HandlerErrorGrpc(err)
	}

	return listProd, nil
}

func (repo *ProductRepository) ProductUploadImages(tenantID string, schema *schemas.ProductUploadSchema, productID int64, validationData *schemas.ProductValidateImage, ctx context.Context) error {
	type imageResult struct {
		primary     *string
		secondaries []string
		filesNames  []string
		uuidBases   []string
		err         error
	}

	var (
		result imageResult
		wg     sync.WaitGroup
		mu     sync.Mutex
	)
	if schema.PrimaryImage != nil {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
		defer wg.Done()
		fileNames, uuidGen, err := utils.SaveTenantImages(tenantID, file, 200, 500)

		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.err = err
			return
		}
		result.primary = &uuidGen
		result.filesNames = append(result.filesNames, fileNames...)
		result.uuidBases = append(result.uuidBases, uuidGen)
		}(schema.PrimaryImage)
	}

	for _, fileImg := range schema.SecondaryImages {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			fileNames, uuidGen, err := utils.SaveTenantImages(tenantID, file, 200, 500)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.err = err
				return
			}
			result.secondaries = append(result.secondaries, uuidGen)
			result.filesNames = append(result.filesNames, fileNames...)
			result.uuidBases = append(result.uuidBases, uuidGen)
		}(fileImg)
	}

	wg.Wait()

	if result.err != nil {
		for _, uuidBase := range result.uuidBases {
			err := utils.DeleteTenantImages(tenantID, uuidBase, 200, 500)
			if err != nil {
				log.Error().Err(err).Msg("Error al eliminar imágenes tras fallo en guardado")
			}
		}
		return schemas.ErrorResponse(500, "Error al guardar imagenes", result.err)
	}

	ctxt, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	md := metadata.Pairs("x-tenant-identifier", tenantID)

	outCtx := metadata.NewOutgoingContext(ctxt, md)

	req := &pb.SaveImageRequest{
		ProdId:          productID,
		PrimaryImage:    result.primary,
		SecondaryImages: result.secondaries,
		KeepSecondaries: validationData.SecondaryImage.KeepUUIDs,
		RemoveSecondaries: validationData.SecondaryImage.RemoveUUIDs,
	}

	_, err := repo.Client.SaveUrlImage(outCtx, req)
	if err != nil {
		return schemas.HandlerErrorGrpc(err)
	}

	for _, uuidBase := range validationData.SecondaryImage.KeepUUIDs {
		err := utils.DeleteTenantImages(tenantID, uuidBase, 200, 500)
		if err != nil {
			log.Error().Err(err).Msg("Error al eliminar imágenes tras guardado")
		}
	}

	return nil
}
