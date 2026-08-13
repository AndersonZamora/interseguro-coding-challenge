// Package config centraliza la carga de configuracion desde variables de
// entorno para la API en Go.
package config

import (
	"fmt"
	"os"
	"time"
)

// Config agrupa toda la configuracion necesaria para levantar el servidor y
// sus dependencias salientes.
type Config struct {
	Port           string
	JWTSecret      string
	NodeAPIURL     string
	NodeAPITimeout time.Duration
	AllowedOrigin  string
}

// Load lee la configuracion desde variables de entorno, aplicando valores
// por defecto razonables para desarrollo local. JWTSecret y NodeAPIURL son
// obligatorios: sin ellos el servicio no puede autenticar ni comunicarse con
// la API de Node, asi que Load falla rapido en vez de arrancar en un estado
// inconsistente.
func Load() (Config, error) {
	cfg := Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		NodeAPIURL:     os.Getenv("NODE_API_URL"),
		NodeAPITimeout: 5 * time.Second,
		AllowedOrigin:  getEnv("ALLOWED_ORIGIN", "*"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es obligatorio")
	}
	if cfg.NodeAPIURL == "" {
		return Config{}, fmt.Errorf("NODE_API_URL es obligatorio")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
