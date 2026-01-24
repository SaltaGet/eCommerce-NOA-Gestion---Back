package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (svc *MPService) MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string, tokenMP *string, baseUrl string, ctx context.Context) (string, error) {
	var listIDs []int64
	for _, item := range data.Items {
		listIDs = append(listIDs, item.ProductID)
	}

	products, err := svc.RepoProd.ProductsValidate(listIDs, tenantIdentifier, ctx)
	if err != nil {
		return "", err
	}

	for _, item := range data.Items {
		for _, prod := range products.Products {
			if item.ProductID == prod.Id {
				if item.UnitPrice != prod.Price {
					return "", schemas.ErrorResponse(400, fmt.Sprintf("El precio del producto %s ha cambiado. Precio actual: %.2f", item.Name, prod.Price), fmt.Errorf("el precio del producto %s ha cambiado. Precio actual: %.2f", item.Name, prod.Price))
				}
				if float64(item.Quantity) > prod.Stock {
					return "", schemas.ErrorResponse(400, fmt.Sprintf("No hay stock suficiente del producto %s. Stock actual: %.2f", item.Name, prod.Stock), fmt.Errorf("no hay stock suficiente del producto %s. Stock actual: %.2f", item.Name, prod.Stock))
				}
				break
			}
		}
	}

	return svc.Repo.MPGenerateLink(data, tenantIdentifier, tokenMP, baseUrl)
}

func (svc *MPService) MPStatePay(data *schemas.WebhookPayload, tenant *schemas.TenantResponseSetting, tokenMP *string, ctx context.Context, ct *fiber.Ctx) error {
	dataPay, err := GetPayInfo(data, tenant, tokenMP)
	if err != nil {
		return err
	}

	if dataPay.Status == "approved" {
		name := dataPay.AditionalInfo.Payer.FirstName + " " + dataPay.AditionalInfo.Payer.LastName
		utils.SendEmailSell(dataPay.Metadata.Email, name, "Pago recibido", *tenant, *dataPay, ct)
	}

	return svc.Repo.MPStatePay(dataPay, tenant, ctx)
}

func (svc *MPService) MPInfoPay(tenantID schemas.TenantResponseSetting, reference string) (*schemas.PaymentSummary, error) {
	accessToken := *tenantID.TokenMP
	url := "https://api.mercadopago.com/v1/payments/search?external_reference=" + reference

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al crear petición", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al contactar Mercado Pago", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, schemas.ErrorResponse(resp.StatusCode, "error al leer respuesta", err)
	}

	summ, err := MapToPaymentSummary(body, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al parsear respuesta", err)
	}

	if resp.StatusCode != 200 {
		return nil, schemas.ErrorResponse(resp.StatusCode, "error al intentar acceder a los datos del pago", fmt.Errorf("%v", summ))
	}

	return summ, nil
}

func GetPayInfo(data *schemas.WebhookPayload, tenant *schemas.TenantResponseSetting, tokenMP *string) (*schemas.DataInfoPay, error) {
	accessToken := *tokenMP
	url := "https://api.mercadopago.com/v1/payments/" + data.Data.ID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al crear petición", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al contactar Mercado Pago", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, schemas.ErrorResponse(resp.StatusCode, "error al leer respuesta", err)
	}

	var dataStruct schemas.DataInfoPay
	err = json.Unmarshal(body, &dataStruct)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al intentar acceder a los datos del pago", err)
	}

	if resp.StatusCode != 200 {
		return nil, schemas.ErrorResponse(resp.StatusCode, fmt.Sprintf("%s", dataStruct.Message), fmt.Errorf("%v", dataStruct))
	}

	if dataStruct.Metadata.TenantIdentifier == "" {
		return nil, schemas.ErrorResponse(500, "error al intentar acceder a los datos del pago", fmt.Errorf("%v", dataStruct))
	}

	return &dataStruct, nil
}

func MapToPaymentSummary(jsonBody []byte, tenant schemas.TenantResponseSetting) (*schemas.PaymentSummary, error) {
	// 1. Unmarshal a un mapa genérico para extraer los datos de "results"
	// Nota: El JSON que pasaste tiene una estructura {"results": [...]}
	var rawResponse struct {
		Results []map[string]interface{} `json:"results"`
	}

	if err := json.Unmarshal(jsonBody, &rawResponse); err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	if len(rawResponse.Results) == 0 {
		return nil, fmt.Errorf("no results found in payment response")
	}

	// Tomamos el primer resultado
	res := rawResponse.Results[0]

	// 2. Extraer datos con aserciones de tipo seguras
	summary := &schemas.PaymentSummary{
		ID:                int64(res["id"].(float64)),
		Status:            res["status"].(string),
		StatusDetail:      res["status_detail"].(string),
		ExternalReference: res["external_reference"].(string),
		Amount:            res["transaction_amount"].(float64),
		Currency:          res["currency_id"].(string),
		PaymentMethodID:   res["payment_method_id"].(string),
		PaymentType:       res["payment_type_id"].(string),
	}

	// 3. Parsear Fecha
	if dateStr, ok := res["date_approved"].(string); ok {
		t, _ := time.Parse(time.RFC3339, dateStr)
		summary.DateApproved = t
	}

	if dateStr, ok := res["date_created"].(string); ok {
		t, _ := time.Parse(time.RFC3339, dateStr)
		summary.DateCreated = t
	}

	// 4. Mapear Payer (Email y Nombre)
	if payer, ok := res["payer"].(map[string]interface{}); ok {
		summary.PayerEmail, _ = payer["email"].(string)
	}

	if addInfo, ok := res["additional_info"].(map[string]interface{}); ok {
		if p, ok := addInfo["payer"].(map[string]interface{}); ok {
			firstName, _ := p["first_name"].(string)
			lastName, _ := p["last_name"].(string)
			summary.PayerName = fmt.Sprintf("%s %s", firstName, lastName)
		}

		// 5. Mapear Items
		if itemsRaw, ok := addInfo["items"].([]interface{}); ok {
			for _, it := range itemsRaw {
				itemMap := it.(map[string]interface{})
				qty, _ := strconv.Atoi(itemMap["quantity"].(string))
				price, _ := strconv.ParseFloat(itemMap["unit_price"].(string), 64)

				summary.Items = append(summary.Items, schemas.Item{
					Title:     itemMap["title"].(string),
					Quantity:  qty,
					UnitPrice: price,
				})
			}
		}
	}

	// 6. Tarjeta (Last Four)
	if card, ok := res["card"].(map[string]interface{}); ok {
		summary.CardLastFour, _ = card["last_four_digits"].(string)
	}

	// 7. Datos del Proveedor (Tenant)
	summary.TenantIdentifier = tenant.Identifier
	summary.TenantEmail = tenant.Email
	var phone string
	if tenant.SettingTenant.Phone != nil {
		phone = *tenant.SettingTenant.Phone
	} else {
		phone = tenant.Phone
	}
	summary.TenantPhone = phone

	return summary, nil
}
 