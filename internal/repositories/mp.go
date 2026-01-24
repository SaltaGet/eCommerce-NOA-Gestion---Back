package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	"google.golang.org/grpc/metadata"
)

func (r *MercadoPagoRepository) MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string, tokenMP *string, baseUrl string) (string, error) {
	accessToken := *tokenMP
	cfg, err := config.New(accessToken)
	if err != nil {
		return "", err
	}

	client := preference.NewClient(cfg)
	UUID := uuid.New().String()
	expiration := time.Now().Add(time.Minute * 30)

	baseUrlDom := strings.Replace(baseUrl, "ecommerce.", "", 1)

	request := preference.Request{

		Payer: &preference.PayerRequest{
			Name:    data.Client.Name,
			Surname: data.Client.Surname,
			Email:   data.Client.Email,
		},
		Metadata: map[string]any{"tenant_identifier": tenantIdentifier, "email": data.Client.Email},
		BackURLs: &preference.BackURLsRequest{
			Success: fmt.Sprintf("%s/store/%s/external_reference/success?uuid=%s", baseUrlDom, tenantIdentifier, UUID),
			Pending: fmt.Sprintf("%s/store/%s/external_reference/pending?uuid=%s", baseUrlDom, tenantIdentifier, UUID),
			Failure: fmt.Sprintf("%s/store/%s/external_reference/failure?uuid=%s", baseUrlDom, tenantIdentifier, UUID),
		},
		ExternalReference: UUID,
		DateOfExpiration:  &expiration,
	}

	for _, item := range data.Items {
		request.Items = append(request.Items, preference.ItemRequest{
			ID:         fmt.Sprintf("%d", item.ProductID),
			Title:      item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			CurrencyID: "ARS",
		})
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		return "", schemas.ErrorResponse(500, "error al genenrar el link de pago", err)
	}

	return resource.InitPoint, nil
}

func (r *MercadoPagoRepository) MPStatePay(data *schemas.DataInfoPay, tenant *schemas.TenantResponseSetting, ctx context.Context) error {
	ctxt, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	md := metadata.Pairs("x-tenant-identifier", tenant.Identifier)
	outCtx := metadata.NewOutgoingContext(ctxt, md)

	dataInfoPay := &pb.DataInfoPay{
		DateApproved:      data.DateApproved,
		DateCreated:       data.DateCreated,
		Status:            data.Status,
		Id:                fmt.Sprintf("%.0f", data.Id),
		StatusDetail:      data.StatusDetail,
		TransactionAmount: data.TransactionAmount,
		TransactionDetails: &pb.TransactionDetails{
			NetReceivedAmount: data.TransactionDetails.NetReceivedAmount,
			TotalPaidAmount:   data.TransactionDetails.TotalPaidAmount,
		},
		Payer: &pb.PayerInfo{
			FirstName: data.AditionalInfo.Payer.FirstName,
			LastName:  data.AditionalInfo.Payer.LastName,
			Email:     data.AditionalInfo.Payer.Email,
		},
		OperationType:     data.OperationType,
		Message:           data.Message,
		ExternalReference: data.ExternalReference,
	}

	dataInfoPay.AdditionalInfo = &pb.AdditionalInfo{}
	dataInfoPay.PaymentMethod = &pb.PaymentMethod{
		Type: data.PayMethod.Type,
	}

	for _, dt := range data.AditionalInfo.Items {
		dataInfoPay.AdditionalInfo.Items = append(dataInfoPay.AdditionalInfo.Items, &pb.ItemsCartPay{
			Id:        dt.ProductID,
			Title:     dt.Title,
			Quantity:  dt.Quantity,
			UnitPrice: dt.UnitPrice,
		})
	}

	dataInfoPay.AdditionalInfo.Payer = &pb.PayerInfo{
		FirstName: data.AditionalInfo.Payer.FirstName,
		LastName:  data.AditionalInfo.Payer.LastName,
		Email:     data.AditionalInfo.Payer.Email,
	}

	_, err := r.Client.SyncPurchasePayment(outCtx, dataInfoPay)
	if err != nil {
		return schemas.HandlerErrorGrpc(err)
	}

	return nil
}
