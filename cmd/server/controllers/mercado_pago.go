package controllers

import (
	"fmt"
	"os"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// MPGenerateLink godoc
//
//	@Summary		MPGenerateLink
//	@Description	Obtener link de pago
//	@Tags			Mercado Pago
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
// @Param		shopping_cart		body		schemas.ShoppingCart	true	"Body"
//	@Success		200			{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/mp/generate_link [post]
func (m *MPController) MPGenerateLink(c *fiber.Ctx) error {
	var shoppingCart schemas.ShoppingCart
	if err := c.BodyParser(&shoppingCart); err != nil {
		return err
	}
	if shoppingCart.Validate() != nil {
		return shoppingCart.Validate()
	}

	tenantIdentifier := c.Locals("tenant_identifier").(string)

	link, err := m.MPService.MPGenerateLink(&shoppingCart, tenantIdentifier)
	if err != nil {
		return err
	}
	return c.JSON(link)
}

// MPStatePay godoc
//
//	@Summary		MPStatePay
//	@Description	Obtener link de pago
//	@Tags			Mercado Pago
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/mp/state_pay [post]
func (m *MPController) MPStatePay(c *fiber.Ctx) error {
	body := c.Body()

	fmt.Println(string(body))

	// {"action":"payment.created","api_version":"v1","data":{"id":"140860642477"},"date_created":"2026-01-10T22:07:06Z","id":127977043877,"live_mode":true,"type":"payment","user_id":"380297435"}
	// paymentID := "127977043877"
	paymentID := "140860642477"
	accessToken := os.Getenv("ACCESS_TOKEN_MP")
	url := "https://api.mercadopago.com/v1/payments/" + paymentID

	agent := fiber.Get(url)
	agent.Set("Authorization", "Bearer "+accessToken)
	agent.Set("Content-Type", "application/json")

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(500).JSON(fiber.Map{"error": errs[0].Error()})
	}

	// campos a tener en cuenta
	// additional_info, date_approved, date_created, id, status, status_detail, transaction_amount, transaction_details

	return c.Status(statusCode).Send(body)
}


// {
//    "accounts_info":null,
//    "acquirer_reconciliation":[
      
//    ],
//    "additional_info":{
//       "items":[
//          {
//             "id":"1",
//             "quantity":"1",
//             "title":"zopeda gigante de madera",
//             "unit_price":"100.0"
//          }
//       ],
//       "payer":{
//          "first_name":"Gaston Gonzalez"
//       },
//       "tracking_id":"platform:v1-whitelabel,so:ALL,type:N/A,security:none"
//    },
//    "authorization_code":null,
//    "binary_mode":false,
//    "brand_id":null,
//    "build_version":"3.137.0-rc-1",
//    "call_for_authorize_id":null,
//    "captured":true,
//    "card":{
      
//    },
//    "charges_details":[
//       {
//          "accounts":{
//             "from":"collector",
//             "to":"mp"
//          },
//          "amounts":{
//             "original":4.1,
//             "refunded":0
//          },
//          "base_amount":100,
//          "client_id":0,
//          "date_created":"2026-01-10T18:07:06.000-04:00",
//          "external_charge_id":"01KEMZ5RZK79A5CGC8FMVD857K",
//          "id":"140860642477-001",
//          "last_updated":"2026-01-10T18:07:06.000-04:00",
//          "metadata":{
//             "reason":"",
//             "source":"proc-svc-charges",
//             "source_detail":"processing_fee_charge"
//          },
//          "name":"mercadopago_fee",
//          "rate":4.1,
//          "refund_charges":[
            
//          ],
//          "reserve_id":null,
//          "type":"fee",
//          "update_charges":[
            
//          ]
//       }
//    ],
//    "charges_execution_info":{
//       "internal_execution":{
//          "date":"2026-01-10T18:07:06.240-04:00",
//          "execution_id":"01KEMZ5RYDXJ6QMM36TNPKGCYN"
//       }
//    },
//    "collector_id":380297435,
//    "corporation_id":null,
//    "counter_currency":null,
//    "coupon_amount":0,
//    "currency_id":"ARS",
//    "date_approved":"2026-01-10T18:07:06.000-04:00",
//    "date_created":"2026-01-10T18:07:06.000-04:00",
//    "date_last_updated":"2026-01-10T18:07:08.000-04:00",
//    "date_of_expiration":null,
//    "deduction_schema":null,
//    "description":"zopeda gigante de madera",
//    "differential_pricing_id":null,
//    "external_reference":null,
//    "fee_details":[
//       {
//          "amount":4.1,
//          "fee_payer":"collector",
//          "type":"mercadopago_fee"
//       }
//    ],
//    "financing_group":null,
//    "id":140860642477,
//    "installments":1,
//    "integrator_id":null,
//    "issuer_id":"2005",
//    "live_mode":true,
//    "marketplace_owner":null,
//    "merchant_account_id":null,
//    "merchant_number":null,
//    "metadata":{
//       "tenant_identifier":"daniel"
//    },
//    "money_release_date":"2026-01-28T18:07:06.000-04:00",
//    "money_release_schema":null,
//    "money_release_status":"pending",
//    "notification_url":null,
//    "operation_type":"regular_payment",
//    "order":{
//       "id":"37197002228",
//       "type":"mercadopago"
//    },
//    "payer":{
//       "email":"pico_elcuervo_sarmiento@hotmail.com",
//       "entity_type":null,
//       "first_name":null,
//       "id":"155963404",
//       "identification":{
//          "number":"20410532132",
//          "type":"CUIL"
//       },
//       "last_name":null,
//       "operator_id":null,
//       "phone":{
//          "number":null,
//          "extension":null,
//          "area_code":null
//       },
//       "type":null
//    },
//    "payment_method":{
//       "id":"account_money",
//       "issuer_id":"2005",
//       "type":"account_money"
//    },
//    "payment_method_id":"account_money",
//    "payment_type_id":"account_money",
//    "platform_id":null,
//    "point_of_interaction":{
//       "business_info":{
//          "branch":"PX",
//          "sub_unit":"checkout_pro",
//          "unit":"online_payments"
//       },
//       "transaction_data":{
         
//       },
//       "type":"UNSPECIFIED"
//    },
//    "pos_id":null,
//    "processing_mode":"aggregator",
//    "refunds":[
      
//    ],
//    "release_info":null,
//    "shipping_amount":0,
//    "sponsor_id":null,
//    "statement_descriptor":null,
//    "status":"approved",
//    "status_detail":"accredited",
//    "store_id":null,
//    "tags":null,
//    "taxes_amount":0,
//    "transaction_amount":100,
//    "transaction_amount_refunded":0,
//    "transaction_details":{
//       "acquirer_reference":null,
//       "external_resource_url":null,
//       "financial_institution":null,
//       "installment_amount":0,
//       "net_received_amount":95.9,
//       "overpaid_amount":0,
//       "payable_deferral_period":null,
//       "payment_method_reference_id":null,
//       "total_paid_amount":100
//    }
// }