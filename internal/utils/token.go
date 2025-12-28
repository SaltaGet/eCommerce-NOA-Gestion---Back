package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString string) (jwt.Claims, error) {
	clearToken := CleanToken(tokenString)
	token, err := jwt.Parse(clearToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY_VALIDATOR")), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, schemas.ErrorResponse(401, "Token vencido", err)
		}
		return nil, schemas.ErrorResponse(401, "Token inválido", err)
	}

	if !token.Valid {
		return nil, schemas.ErrorResponse(401, "Token inválido", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, schemas.ErrorResponse(500, "No se pudieron obtener los claims", nil)
	}

	return claims, nil
}

func CleanToken(bearerToken string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(bearerToken, prefix) {
		return strings.TrimPrefix(bearerToken, prefix)
	}
	return bearerToken
}