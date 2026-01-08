package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
)

func (s *TenantService) TenantList() ([]schemas.TenantResponse, error) {
	resp, err := s.Repo.TenantList()
	if err != nil {
		return nil, err
	}

	var tenants []schemas.TenantResponse
	for _, t := range resp.Tenants {
		var expiration *time.Time
		if t.Expiration != nil {
			exp := t.Expiration.AsTime()
			expiration = &exp
		}
		tenant := schemas.TenantResponse{
			Identifier: t.Identifier,
			IsActive:   t.IsActive,
			Expiration: expiration,
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (s *TenantService) TenantGet(tenantIdentifier string) (*schemas.TenantResponseSetting, error) {
	resp, err := s.Repo.TenantGet(tenantIdentifier)
	if err != nil {
		return nil, err
	}

	tenant := schemas.TenantResponseSetting{
		ID:         resp.Id,
		Name:       resp.Name,
		Identifier: resp.Identifier,
		Address:    resp.Address,
		Phone:      resp.Phone,
		Email:      resp.Email,
		SettingTenant: schemas.SettingTenant{
			ID:             resp.SettingTenant.Id,
			LogoSmall:      &resp.SettingTenant.Logo,
			LogoBig:        &resp.SettingTenant.Logo,
			FrontPageSmall: &resp.SettingTenant.FrontPage,
			FrontPageBig:   &resp.SettingTenant.FrontPage,
			Title:          &resp.SettingTenant.Title,
			Slogan:         &resp.SettingTenant.Slogan,
			PrimaryColor:   &resp.SettingTenant.PrimaryColor,
			SecondaryColor: &resp.SettingTenant.SecondaryColor,
		},
	}

	return &tenant, nil
}

func (s *TenantService) TenantSaveImage(tenantID string, schema *schemas.TenantUploadSchema, ctx context.Context) error {
	var logoUUID, frontPageUUID string
	var err error
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	if schema.LogoImage == nil {
		logoUUID = ""
	} else {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			_, lUUID, e := utils.SaveTenantImages(tenantID, file, 200, 500)
			mu.Lock()
			defer mu.Unlock()
			if e != nil {
				errs = append(errs, errors.New("Imágen del logo"))
				return
			}

			logoUUID = lUUID

			return
		}(schema.LogoImage)
		// logoUrl, logoUUID, err = utils.SaveTenantImages(tenantID, schema.LogoImage, 200, 500)
		// if err != nil {
		// 	return schemas.ErrorResponse(500, err.Error(), err)
		// }
	}

	if schema.FrontPageImages == nil {
		frontPageUUID = ""
	} else {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			_, frontUUID, e := utils.SaveTenantImages(tenantID, schema.FrontPageImages, 500, 1000)
			mu.Lock()
			defer mu.Unlock()
			if e != nil {
				errs = append(errs, errors.New("Imágen de la portada"))
				return
			}

			frontPageUUID = frontUUID

			return
		}(schema.FrontPageImages)
	}

	wg.Wait()

	req := &pb.TenantRequestImageSetting{
		LogoUuid:      utils.Ternary(logoUUID == "", nil, &logoUUID),
		FrontPageUuid: utils.Ternary(frontPageUUID == "", nil, &frontPageUUID),
		TenantIdentifier: tenantID,
	}

	resp, err := s.Repo.TenantSaveImage(req, ctx)
	if err != nil {
		_ = utils.DeleteTenantImages(tenantID, logoUUID, 200, 500)
		_ = utils.DeleteTenantImages(tenantID, frontPageUUID, 500, 1000)
		return err
	}

	if resp.LogoUuid != nil {
		_ = utils.DeleteTenantImages(tenantID, *resp.LogoUuid, 200, 500)
	}

	if resp.FrontPageUuid != nil {
		_ = utils.DeleteTenantImages(tenantID, *resp.FrontPageUuid, 500, 1000)
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("se produjo un error, por favor intente nuevamente en caso de que alguna de las siguientes imagenes fallara: %w", errors.Join(errs...))
	}

	return nil
}
