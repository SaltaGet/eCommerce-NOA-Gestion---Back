package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

func (r *MercadoPagoRepository) MPGenerateLink(data *schemas.ShoppingCart, tenantIdentifier string) (string, error) {
	accessToken := os.Getenv("ACCESS_TOKEN_MP")
	// accessToken := "APP_USR-3934855919608956-011003-1954eda055464883314304027aec5749-3125334152" // Prueba
	//obtener la key del usuario para realizar pagos

	cfg, err := config.New(accessToken)
	if err != nil {
		return "", err
	}

	client := preference.NewClient(cfg)

	request := preference.Request{
		
		Payer: &preference.PayerRequest{
			Name:  "Gaston Gonzalez",
			Surname: "Gonzalez",
			Email: "gonzalezgastonariel@gmail.com",
		},
		Items: []preference.ItemRequest{
			{
				ID:        "1",
				Title:     "zopeda gigante de madera",
				Quantity:  1,
				UnitPrice: 100.00,
				CurrencyID: "ARS",
			},
		},
		Metadata: map[string]any{"tenant_identifier": "daniel"},
	}

	for _, item := range data.Items {
		request.Items = append(request.Items, preference.ItemRequest{
			ID:        string(item.ProductID),
			Title:     item.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			CurrencyID: "ARS",
		})
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	fmt.Printf("%+#v\n", resource)

	bytes, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		return "", err
	}

	fmt.Println(string(bytes))

	return resource.InitPoint, nil
}
