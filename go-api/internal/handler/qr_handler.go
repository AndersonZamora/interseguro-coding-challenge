// Package handler contiene los handlers HTTP (Fiber) de la API de Go.
package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/interseguro/coding-challenge/go-api/internal/client"
	"github.com/interseguro/coding-challenge/go-api/internal/matrix"
)

// StatsFetcher abstrae la llamada a la API de Node para poder inyectar un
// mock en los tests del handler sin depender de un servidor HTTP real.
type StatsFetcher interface {
	FetchStats(ctx context.Context, q, r [][]float64) (client.StatsResponse, error)
}

// Handler agrupa las dependencias necesarias para atender las rutas de
// matrices.
type Handler struct {
	statsFetcher StatsFetcher
	callTimeout  time.Duration
}

// NewHandler crea un Handler. callTimeout acota cuanto se espera la
// respuesta de la API de Node antes de devolver 502 al cliente.
func NewHandler(statsFetcher StatsFetcher, callTimeout time.Duration) *Handler {
	return &Handler{statsFetcher: statsFetcher, callTimeout: callTimeout}
}

type qrRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type qrResponse struct {
	Q     [][]float64          `json:"q"`
	R     [][]float64          `json:"r"`
	Stats client.StatsResponse `json:"stats"`
}

// PostQR atiende POST /api/v1/qr: valida y factoriza la matriz de entrada,
// reenvia Q y R a la API de Node para obtener estadisticas, y devuelve todo
// combinado.
func (h *Handler) PostQR(c *fiber.Ctx) error {
	var req qrRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body",
		})
	}

	result, err := matrix.Decompose(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), h.callTimeout)
	defer cancel()

	stats, err := h.statsFetcher.FetchStats(ctx, result.Q, result.R)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "node-api unavailable",
		})
	}

	return c.Status(fiber.StatusOK).JSON(qrResponse{
		Q:     result.Q,
		R:     result.R,
		Stats: stats,
	})
}
