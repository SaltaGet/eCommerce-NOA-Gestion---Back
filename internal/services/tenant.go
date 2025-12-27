package services

import (
	"time"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
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