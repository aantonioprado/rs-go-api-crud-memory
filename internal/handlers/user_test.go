package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aantonioprado/rs-go-api-crud-memory/internal/handlers"
	"aantonioprado/rs-go-api-crud-memory/internal/store"
	"aantonioprado/rs-go-api-crud-memory/internal/utils"

	"github.com/go-chi/chi/v5"
)

func setupRouter() *chi.Mux {
	mem := store.NewMemory()
	h := handlers.NewUserHandler(mem)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotFound, "Route not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed for this route")
	})
	return r
}

func TestUsers_Create_List_Get_Update_Delete(t *testing.T) {
	r := setupRouter()

	body := `{"first_name":"Jane","last_name":"Doe","biography":"Biografia válida com >= 20 caracteres."}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201. body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("POST decode response: %v", err)
	}
	if !resp.Success || resp.Message == "" {
		t.Fatalf("POST expected success with message, got: %+v", resp)
	}

	var created struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Biography string `json:"biography"`
	}
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		t.Fatalf("POST decode data: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("POST expected non-empty UUID id")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET all status = %d, want 200. body=%s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET by id status = %d, body=%s", rec.Code, rec.Body.String())
	}

	updBody := `{"first_name":"Janette","last_name":"Doe","biography":"Biografia atualizada com tamanho válido (>= 20)."}`
	req = httptest.NewRequest(http.MethodPut, "/api/users/"+created.ID, bytes.NewBufferString(updBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, body=%s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("DELETE status = %d, body=%s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET after delete status = %d, want 404, body=%s", rec.Code, rec.Body.String())
	}
}

func TestUsers_BadRequest_On_InvalidPayload(t *testing.T) {
	r := setupRouter()

	body := `{"first_name":"A","last_name":"Doe","biography":"curta"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST invalid payload -> status = %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"success":false`) {
		t.Fatalf("POST invalid payload should be error JSON: %s", rec.Body.String())
	}
}

func TestUsers_NotFoundRoute(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/does-not-exist", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	if got := rec.Body.String(); !strings.Contains(got, `"message":"Route not found"`) {
		t.Fatalf("missing Route not found message: %s", got)
	}
}
