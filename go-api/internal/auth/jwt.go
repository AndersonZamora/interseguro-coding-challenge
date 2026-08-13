// Package auth implementa emision y verificacion de JWT compartidos entre
// la API de Go y la API de Node (mismo secreto HMAC via env var JWT_SECRET).
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ErrMissingOrInvalidToken es el error uniforme devuelto (como 401) cuando
// falta el header Authorization o el token no es valido/expiro. El mismo
// mensaje se usa en la API de Node para que el frontend maneje un solo caso.
var ErrMissingOrInvalidToken = errors.New("missing or invalid authorization token")

const serviceTokenTTL = 5 * time.Minute

// IssueServiceToken firma un JWT de corta duracion para la llamada
// service-to-service que la API de Go hace hacia la API de Node, sin pasar
// por el endpoint de login (que existe en Node para clientes humanos/UI).
func IssueServiceToken(secret string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   "go-api-service",
		Issuer:    "interseguro-challenge",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(serviceTokenTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Middleware valida el Bearer token de cada request usando el secreto
// compartido. Devuelve 401 con un mensaje uniforme si el token falta,
// esta mal formado, expiro o la firma no coincide.
func Middleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := extractBearerToken(c.Get(fiber.HeaderAuthorization))
		if raw == "" {
			return unauthorized(c)
		}

		_, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil {
			return unauthorized(c)
		}

		return c.Next()
	}
}

func extractBearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": ErrMissingOrInvalidToken.Error(),
	})
}
