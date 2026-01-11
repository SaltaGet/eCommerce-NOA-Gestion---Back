package ports

import "github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"

type MPRepository interface {
	MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string) (string, error)
}

type MPService interface {
	MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string) (string, error)
}
