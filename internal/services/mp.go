package services

import "github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"

func (svc *MPService) MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string) (string, error) {
	return svc.Repo.MPGenerateLink(data, tenantIdentifier)
}