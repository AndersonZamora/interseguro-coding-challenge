package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/interseguro/coding-challenge/go-api/internal/client"
)

// stubStatsFetcher es un StatsFetcher falso controlado por el test, sin
// necesidad de levantar un servidor HTTP real para la API de Node.
type stubStatsFetcher struct {
	response client.StatsResponse
	err      error
}

func (s stubStatsFetcher) FetchStats(ctx context.Context, q, r [][]float64) (client.StatsResponse, error) {
	return s.response, s.err
}

func newApp(fetcher StatsFetcher) *fiber.App {
	app := fiber.New()
	h := NewHandler(fetcher, 5*time.Second)
	app.Post("/api/v1/qr", h.PostQR)
	return app
}

func TestPostQR_Success(t *testing.T) {
	fetcher := stubStatsFetcher{response: client.StatsResponse{Max: 5, Min: 1, Average: 3, Sum: 12, AnyDiagonal: false}}
	app := newApp(fetcher)

	body, _ := json.Marshal(map[string]any{
		"matrix": [][]float64{{1, 2}, {3, 4}, {5, 6}},
	})
	req := httptest.NewRequest("POST", "/api/v1/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var got qrResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Stats.Max != 5 {
		t.Errorf("stats.Max = %v, want 5", got.Stats.Max)
	}
	if len(got.Q) != 3 || len(got.R) != 2 {
		t.Errorf("unexpected Q/R dimensions: Q=%v R=%v", got.Q, got.R)
	}
}

func TestPostQR_InvalidMatrix(t *testing.T) {
	app := newApp(stubStatsFetcher{})

	body, _ := json.Marshal(map[string]any{
		"matrix": [][]float64{},
	})
	req := httptest.NewRequest("POST", "/api/v1/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestPostQR_NodeAPIUnavailable(t *testing.T) {
	fetcher := stubStatsFetcher{err: errors.New("connection refused")}
	app := newApp(fetcher)

	body, _ := json.Marshal(map[string]any{
		"matrix": [][]float64{{1, 2}, {3, 4}},
	})
	req := httptest.NewRequest("POST", "/api/v1/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
}
