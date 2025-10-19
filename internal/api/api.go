package api

import (
	"aantonioprado/rs-go-api-crud-memory/internal/handlers"
	"aantonioprado/rs-go-api-crud-memory/internal/store"
	"aantonioprado/rs-go-api-crud-memory/internal/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 & time.Second))

	memory := store.NewMemory()
	usersHandler := handlers.NewUserHandler(memory)

	usersHandler.RegisterRoutes(r)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotFound, "Route not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed for this route")
	})

	return r
}
