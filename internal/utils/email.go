package utils

import (
	"fmt"
	"io"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
)

func SendEmailSell(to, name, Subject string, tenant schemas.TenantResponseSetting, data schemas.DataInfoPay, ctx *fiber.Ctx) error {
	m := gomail.NewMessage()
	m.SetHeader("From", tenant.Email)
	m.SetHeader("To", to) // Usar el parámetro 'to'
	m.SetHeader("Subject", Subject)

	colorP := "#164e43"
	colorS := "#06b6d4"
	if tenant.SettingTenant.PrimaryColor != nil && tenant.SettingTenant.SecondaryColor != nil {
		colorP = *tenant.SettingTenant.PrimaryColor
		colorS = *tenant.SettingTenant.SecondaryColor
	}

	// Teléfono seguro
	tenantPhone := tenant.Phone
	if tenant.SettingTenant.Phone != nil {
		tenantPhone = *tenant.SettingTenant.Phone
	}

	// 1. Embeber Logo
	if tenant.SettingTenant.LogoBig != nil {
		path, exists := GetPath(tenant.Identifier, *tenant.SettingTenant.LogoBig+"p500.webp")
		if exists {
			m.Embed(path, gomail.SetHeader(map[string][]string{"Content-ID": {"<logo>"}}))
		}
	}

	// 2. Generar y Embeber Código QR
	// Reemplaza 'http://tudominio.com' por tu variable c.BaseUrl
	qrUrl := fmt.Sprintf("%s/store/%s/info_pay/%s",
		ctx.BaseURL(), data.Metadata.TenantIdentifier, data.ExternalReference)

	qrPng, err := qrcode.Encode(qrUrl, qrcode.Medium, 256)
	if err == nil {
		m.Embed("qrcode.png", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(qrPng)
			return err
		}), gomail.SetHeader(map[string][]string{"Content-ID": {"<qrcode>"}}))
	}

	fullBody := BuildPaymentEmailBody(data, colorP, colorS, tenantPhone, tenant.Address)
	m.SetBody("text/html", fullBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, tenant.Email, *tenant.TokenEmail)
	return d.DialAndSend(m)
}

func BuildPaymentEmailBody(data schemas.DataInfoPay, primaryColor, secondaryColor, tenantPhone, tenantAddress string) string {
	var itemsHtml string
	for _, item := range data.AditionalInfo.Items {
		itemsHtml += fmt.Sprintf(`
            <tr>
                <td style="padding: 15px 12px; border-bottom: 1px solid %s;">
                    <span style="font-weight: bold; color: #333; font-size: 16px;">%s</span><br>
                    <span style="font-size: 14px; color: #777;">Cantidad: %s</span>
                </td>
                <td style="padding: 15px 12px; border-bottom: 1px solid %s; text-align: right; color: %s; font-weight: bold; font-size: 16px;">$%s</td>
            </tr>`, secondaryColor, item.Title, item.Quantity, secondaryColor, primaryColor, item.UnitPrice)
	}

	statusBg := primaryColor
	if data.Status != "approved" {
		statusBg = "#f39c12"
	}

	return fmt.Sprintf(`
    <div style="font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; color: #333; border: 1px solid %s; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 15px rgba(0,0,0,0.08);">
        
        <div style="background-color: %s; padding: 40px 20px; text-align: center; color: white;">
            <img src="cid:logo" alt="Logo" style="max-width: 180px; margin-bottom: 20px; filter: drop-shadow(0px 2px 4px rgba(0,0,0,0.2));">
            <h1 style="margin: 0; font-size: 28px; text-transform: uppercase; letter-spacing: 1px; font-weight: 800;">¡Gracias por tu compra!</h1>
        </div>
        
        <div style="padding: 35px; background-color: white;">
            <div style="display: inline-block; padding: 8px 16px; border-radius: 25px; background-color: %s; color: white; font-size: 14px; font-weight: bold; margin-bottom: 25px; text-transform: uppercase;">
                ESTADO: %s
            </div>

            <p style="font-size: 18px; margin-bottom: 10px;">Hola <strong>%s %s</strong>,</p>
            <p style="color: #555; line-height: 1.6; font-size: 16px;">Confirmamos que hemos recibido tu pago correctamente. Aquí tienes el detalle actualizado de tu pedido:</p>
            
            <table style="width: 100%%; border-collapse: collapse; margin: 25px 0;">
                <thead>
                    <tr style="background-color: %s; color: #333;">
                        <th style="padding: 15px 12px; text-align: left; border-radius: 8px 0 0 8px; font-size: 15px;">Producto</th>
                        <th style="padding: 15px 12px; text-align: right; border-radius: 0 8px 8px 0; font-size: 15px;">Subtotal</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
            </table>

            <div style="background-color: %s; padding: 25px; border-radius: 12px; border: 1px solid %s; margin-bottom: 30px;">
                <table style="width: 100%%;">
                    <tr>
                        <td style="font-size: 15px; color: #555; padding-bottom: 8px;">Método de Pago:</td>
                        <td style="text-align: right; font-size: 15px; font-weight: bold; padding-bottom: 8px;">%s</td>
                    </tr>
                    <tr>
                        <td style="font-size: 15px; color: #555; padding-bottom: 8px;">ID de Transacción:</td>
                        <td style="text-align: right; font-size: 15px; font-weight: bold; padding-bottom: 8px;">#%0.0f</td>
                    </tr>
                    <tr>
                        <td style="padding-top: 15px; font-size: 20px; font-weight: bold; color: %s;">Total Pagado:</td>
                        <td style="padding-top: 15px; text-align: right; font-size: 26px; font-weight: bold; color: %s;">$%0.2f</td>
                    </tr>
                </table>
            </div>

            <div style="text-align: center; padding: 25px; border: 2px dashed %s; border-radius: 12px; background-color: #fafafa;">
                <p style="margin: 0; font-size: 14px; color: #666; font-weight: bold;">REFERENCIA DE PEDIDO</p>
                <p style="margin: 5px 0 15px 0; font-size: 19px; color: #333; font-weight: 800; word-break: break-all;">%s</p>
                <div style="background-color: #ffffff; padding: 10px; display: inline-block; border-radius: 8px; border: 1px solid #eee;">
                    <img src="cid:qrcode" style="width: 150px; height: 150px; display: block;">
                </div>
            </div>
            
            <div style="margin-top: 35px; text-align: center; border-top: 1px solid #eee; padding-top: 30px;">
                <p style="font-size: 17px; color: #333; font-weight: bold; margin-bottom: 10px;">¿Necesitas ayuda con tu entrega?</p>
                <p style="font-size: 16px; color: #555; margin-bottom: 15px; line-height: 1.5;">
                    Por cualquier duda o para coordinar el retiro, por favor <strong>ponte en contacto con nosotros al siguiente número:</strong>
                </p>
                <p style="font-size: 22px; font-weight: bold; color: %s; margin: 10px 0;">
                    📞 %s
                </p>
                
                <div style="margin-top: 20px; padding: 15px; background-color: #f9f9f9; border-radius: 8px; display: inline-block; width: 90%%;">
                    <p style="margin: 0; font-size: 14px; color: #777; text-transform: uppercase; font-weight: bold;">📍 Nuestra Dirección:</p>
                    <p style="margin: 5px 0 0 0; font-size: 16px; color: #444; font-weight: 500;">%s</p>
                </div>
            </div>
        </div>

        <div style="background-color: #f4f4f4; padding: 30px; text-align: center; font-size: 14px; color: #777; border-top: 1px solid #eee;">
            <p style="margin: 0; color: #444; font-weight: bold; font-size: 16px;">%s</p>
            <p style="margin: 8px 0 0 0;">Has recibido este correo por tu compra realizada en nuestra tienda oficial.</p>
            <p style="margin: 15px 0 0 0; font-size: 12px; color: #aaa;">© 2026 - Sistema de Notificaciones de Pago</p>
        </div>
    </div>`,
		secondaryColor,
		primaryColor,                // Header background
		statusBg, data.StatusDetail, // Badge
		data.AditionalInfo.Payer.FirstName, data.AditionalInfo.Payer.LastName,
		secondaryColor, // Table header background
		itemsHtml,
		secondaryColor, secondaryColor, // Summary box
		data.PayMethod.Type, data.Id,
		primaryColor, primaryColor, data.TransactionAmount, // Total text
		primaryColor,              // Dashed border color
		data.ExternalReference,    // Orden externa
		primaryColor, tenantPhone, // Sección de contacto con mensaje directo
		tenantAddress,                  // Nueva sección de dirección
		data.Metadata.TenantIdentifier) // Nombre del comercio
}
