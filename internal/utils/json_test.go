package utils_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"aantonioprado/rs-go-api-crud-memory/internal/utils"
)

func TestDecodeJSON_Empty(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	var v struct{}
	if err := utils.DecodeJSON(req, &v); err == nil {
		t.Fatalf("expected error for empty body")
	}
}

func TestWriteJSON_And_WriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	utils.WriteJSON(rr, 200, "ok", map[string]int{"n": 1})
	if rr.Code != 200 {
		t.Fatalf("WriteJSON wrong status: %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"success":true`) {
		t.Fatalf("WriteJSON missing success=true: %s", body)
	}
	if !strings.Contains(body, `"message":"ok"`) {
		t.Fatalf("WriteJSON missing message: %s", body)
	}
	if !strings.Contains(body, `"n":1`) {
		t.Fatalf("WriteJSON missing data payload: %s", body)
	}

	rr = httptest.NewRecorder()
	utils.WriteError(rr, 404, "Route not found")
	if rr.Code != 404 {
		t.Fatalf("WriteError wrong status: %d", rr.Code)
	}
	body = rr.Body.String()
	if !strings.Contains(body, `"success":false`) {
		t.Fatalf("WriteError missing success=false: %s", body)
	}
	if !strings.Contains(body, `"message":"Route not found"`) {
		t.Fatalf("WriteError missing message: %s", body)
	}
	if !strings.Contains(body, `"error":"Not Found"`) {
		t.Fatalf("WriteError missing error status text: %s", body)
	}
}
