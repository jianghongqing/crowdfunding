package api

import (
	"net/http/httptest"
	"testing"
)

func TestParsePagination(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/campaigns", nil)
		limit, offset, err := parsePagination(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if limit != 20 || offset != 0 {
			t.Fatalf("unexpected pagination: limit=%d offset=%d", limit, offset)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/campaigns?limit=50&offset=10", nil)
		limit, offset, err := parsePagination(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if limit != 50 || offset != 10 {
			t.Fatalf("unexpected pagination: limit=%d offset=%d", limit, offset)
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/campaigns?limit=101", nil)
		if _, _, err := parsePagination(req); err == nil {
			t.Fatal("expected error for invalid limit")
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/campaigns?offset=-1", nil)
		if _, _, err := parsePagination(req); err == nil {
			t.Fatal("expected error for invalid offset")
		}
	})
}
