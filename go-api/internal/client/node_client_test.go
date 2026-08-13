package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret"

func TestFetchStats_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Errorf("expected Authorization header with Bearer token, got %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/api/v1/stats" {
			t.Errorf("expected path /api/v1/stats, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatsResponse{
			Max: 10, Min: 1, Average: 5.5, Sum: 55, AnyDiagonal: true,
		})
	}))
	defer server.Close()

	c := NewNodeStatsClient(server.URL, testSecret, 5*time.Second)
	stats, err := c.FetchStats(context.Background(), [][]float64{{1, 2}}, [][]float64{{3, 4}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Max != 10 || stats.Min != 1 || !stats.AnyDiagonal {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestFetchStats_NodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"boom"}`))
	}))
	defer server.Close()

	c := NewNodeStatsClient(server.URL, testSecret, 5*time.Second)
	_, err := c.FetchStats(context.Background(), [][]float64{{1}}, [][]float64{{1}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchStats_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := NewNodeStatsClient(server.URL, testSecret, 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := c.FetchStats(ctx, [][]float64{{1}}, [][]float64{{1}})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
