package ports

import (
	"context"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

type MPRepository interface {
	MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string, tokenMP *string, baseUrl string) (string, error)
	MPStatePay(data *schemas.DataInfoPay, tenanr *schemas.TenantResponseSetting, ctx context.Context) error
}

type MPService interface {
	MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string, tokenMP *string, baseUrl string, ctx context.Context) (string, error)
	MPStatePay(data *schemas.WebhookPayload, tenant *schemas.TenantResponseSetting, tokenMP *string, ctx context.Context, ct *fiber.Ctx) error
	MPInfoPay(tenantID schemas.TenantResponseSetting, reference string) (*schemas.PaymentSummary, error)
}
