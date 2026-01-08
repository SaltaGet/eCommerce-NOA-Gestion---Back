package utils

import (
	"fmt"
	"image"
	_ "image/gif"  // Registra el decodificador GIF
	_ "image/jpeg" // Registra el decodificador JPEG
	_ "image/png"  // Registra el decodificador PNG
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func GenerateUrl(ctx *fiber.Ctx, tenantID string, filename *string, size string) *string {
	if filename == nil {
		return nil
	}

	filenameSize := fmt.Sprintf("%s%s.webp", *filename, size)
	url := fmt.Sprintf("%s/ecommerce/%s/api/v1/image/get/%s", ctx.BaseURL(), tenantID, filenameSize)
	return &url
}

func GetPath(tenantID, filename string) (string, bool) {
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join("media", tenantID, safeFilename)

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return "", false
	}

	return filePath, true
}

func IsWebP(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Leemos los primeros 512 bytes para detectar el tipo
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)
	return contentType == "image/webp"
}

func SaveTenantImages(tenantID string, file *multipart.FileHeader, sizeSmall, sizeBig uint) ([]string, string, error) {
	// 1. Validar tamaño (2MB)
	if file.Size > 2*1024*1024 {
		return nil, "", fmt.Errorf("el archivo excede el límite de 2MB")
	}

	// 2. Generar un único UUID v7 para ambas imágenes
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", fmt.Errorf("error generando UUID: %w", err)
	}
	baseUUID := id.String()

	// 3. Crear carpetas
	targetDir := filepath.Join("media", tenantID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, "", err
	}

	// 4. Decodificar imagen original
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("formato no soportado: %w", err)
	}

	// 5. Configuración de tamaños y sufijos
	variants := []struct {
		suffix string
		pixels uint
	}{
		{fmt.Sprintf("p%d", sizeSmall), sizeSmall},
		{fmt.Sprintf("p%d", sizeBig), sizeBig},
	}

	var savedPaths []string

	// 6. Procesar y guardar cada variante
	for _, v := range variants {
		// Redimensionar manteniendo proporción
		newImg := resize.Thumbnail(v.pixels, v.pixels, img, resize.Lanczos3)

		// Construir nombre: UUID + sufijo + .webp
		fileName := fmt.Sprintf("%s%s.webp", baseUUID, v.suffix)
		destPath := filepath.Join(targetDir, fileName)

		// Crear archivo físico
		f, err := os.Create(destPath)
		if err != nil {
			return nil, "", err
		}

		// Codificar a WebP (Calidad 80)
		err = webp.Encode(f, newImg, &webp.Options{Lossless: false, Quality: 80})
		f.Close()

		if err != nil {
			return nil, "", err
		}

		savedPaths = append(savedPaths, destPath)
	}

	return savedPaths, baseUUID, nil
}

func DeleteTenantImages(tenantID string, baseUUID string, sizeSmall, sizeBig uint) error {
	if baseUUID == "" {
		return nil
	}
	// 1. Definir los sufijos que manejamos en la subida
	variants := []string{fmt.Sprintf("p%d", sizeSmall), fmt.Sprintf("p%d", sizeBig)}

	targetDir := filepath.Join("media", tenantID)
	var errorsList []error

	// 2. Iterar sobre las variantes y eliminarlas
	for _, suffix := range variants {
		fileName := fmt.Sprintf("%s%s.webp", baseUUID, suffix)
		fullPath := filepath.Join(targetDir, fileName)

		// Verificamos si el archivo existe antes de intentar borrar
		if _, err := os.Stat(fullPath); err == nil {
			err := os.Remove(fullPath)
			if err != nil {
				errorsList = append(errorsList, fmt.Errorf("error al eliminar %s: %w", fileName, err))
			}
		}
	}

	// 3. Si hubo errores al borrar algún archivo, los reportamos
	if len(errorsList) > 0 {
		return fmt.Errorf("problemas al eliminar variantes: %v", errorsList)
	}

	return nil
}
