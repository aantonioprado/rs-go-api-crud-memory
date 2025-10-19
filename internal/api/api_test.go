package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"aantonioprado/rs-go-api-crud-memory/internal/api"
)

func TestAPI_NotFound(t *testing.T) {
	r := api.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}
