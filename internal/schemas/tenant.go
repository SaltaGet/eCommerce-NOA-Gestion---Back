package schemas

import "time"

type TenantResponse struct {
	Identifier string    `json:"identifier"`
	IsActive   bool      `json:"is_active"`
	Expiration *time.Time `json:"expiration"`
}
