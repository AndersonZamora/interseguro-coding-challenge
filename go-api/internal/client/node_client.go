// Package client implementa la llamada HTTP saliente desde la API de Go
// hacia la API de Node para obtener estadisticas de las matrices Q y R.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/interseguro/coding-challenge/go-api/internal/auth"
)

// StatsResponse refleja el contrato de respuesta de POST /api/v1/stats en la
// API de Node.
type StatsResponse struct {
	Max         float64 `json:"max"`
	Min         float64 `json:"min"`
	Average     float64 `json:"average"`
	Sum         float64 `json:"sum"`
	AnyDiagonal bool    `json:"anyDiagonal"`
}

type statsRequest struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// NodeStatsClient llama al endpoint de estadisticas de la API de Node,
// autenticandose con un JWT de servicio firmado con el secreto compartido.
type NodeStatsClient struct {
	baseURL    string
	jwtSecret  string
	httpClient *http.Client
}

// NewNodeStatsClient crea un cliente contra baseURL (ej. http://node-api:3000)
// usando jwtSecret para firmar el token de servicio en cada llamada.
func NewNodeStatsClient(baseURL, jwtSecret string, timeout time.Duration) *NodeStatsClient {
	return &NodeStatsClient{
		baseURL:   baseURL,
		jwtSecret: jwtSecret,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchStats envia Q y R a la API de Node y devuelve las estadisticas
// calculadas. ctx debe llevar un deadline razonable (recomendado 5s); si
// Node no responde a tiempo o devuelve un error, FetchStats devuelve un
// error describiendo la causa para que el handler pueda mapearlo a 502.
func (c *NodeStatsClient) FetchStats(ctx context.Context, q, r [][]float64) (StatsResponse, error) {
	body, err := json.Marshal(statsRequest{Q: q, R: r})
	if err != nil {
		return StatsResponse{}, fmt.Errorf("marshal stats request: %w", err)
	}

	token, err := auth.IssueServiceToken(c.jwtSecret)
	if err != nil {
		return StatsResponse{}, fmt.Errorf("issue service token: %w", err)
	}

	url := c.baseURL + "/api/v1/stats"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return StatsResponse{}, fmt.Errorf("build stats request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return StatsResponse{}, fmt.Errorf("call node-api: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatsResponse{}, fmt.Errorf("read node-api response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return StatsResponse{}, fmt.Errorf("node-api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var stats StatsResponse
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return StatsResponse{}, fmt.Errorf("decode node-api response: %w", err)
	}

	return stats, nil
}
