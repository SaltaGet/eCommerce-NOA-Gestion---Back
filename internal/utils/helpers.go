package utils

import "github.com/google/uuid"

func Ternary[T any](cond bool, a, b T) T {
    if cond {
        return a
    }
    return b
}

func IsValidUUIDv7(id string) bool {
	// 1. Intentar parsear el string
	parsed, err := uuid.Parse(id)
	if err != nil {
		// No es un UUID válido en absoluto
		return false
	}

	// 2. Verificar si la versión es la 7
	return parsed.Version() == 7
}