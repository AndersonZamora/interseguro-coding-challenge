// Command api levanta el servidor HTTP de la API de Go: recibe una matriz,
// calcula su factorizacion QR y reenvia el resultado a la API de Node para
// obtener estadisticas.
package main

import (
	"log"

	"github.com/interseguro/coding-challenge/go-api/internal/client"
	"github.com/interseguro/coding-challenge/go-api/internal/config"
	"github.com/interseguro/coding-challenge/go-api/internal/httpserver"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	statsFetcher := client.NewNodeStatsClient(cfg.NodeAPIURL, cfg.JWTSecret, cfg.NodeAPITimeout)
	app := httpserver.New(cfg, statsFetcher)

	log.Printf("go-api listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
