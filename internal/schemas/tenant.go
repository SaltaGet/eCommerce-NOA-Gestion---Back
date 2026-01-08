package schemas

import (
	"mime/multipart"
	"time"
)

type TenantResponse struct {
	Identifier string     `json:"identifier"`
	IsActive   bool       `json:"is_active"`
	Expiration *time.Time `json:"expiration"`
}

type SettingTenant struct {
	ID             int64   `json:"id"`
	LogoSmall      *string `json:"logo_small"`
	LogoBig        *string `json:"logo_big"`
	FrontPageSmall *string `json:"front_page_small"`
	FrontPageBig   *string `json:"front_page_big"`
	Title          *string `json:"title"`
	Slogan         *string `json:"slogan"`
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
}

type TenantResponseSetting struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Identifier    string        `json:"identifier"`
	Address       string        `json:"address"`
	Phone         string        `json:"phone"`
	Email         string        `json:"email"`
	SettingTenant SettingTenant `json:"setting_tenant"`
}

type TenantUploadSchema struct {
	LogoImage       *multipart.FileHeader `form:"logo_image"`
	FrontPageImages *multipart.FileHeader `form:"font_page_image"`
}
