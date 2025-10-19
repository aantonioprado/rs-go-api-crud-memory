package handlers

import (
	"aantonioprado/rs-go-api-crud-memory/internal/models"
	"aantonioprado/rs-go-api-crud-memory/internal/store"
	"aantonioprado/rs-go-api-crud-memory/internal/utils"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	db *store.Memory
}

func NewUserHandler(db *store.Memory) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", u.createUser)
		r.Get("/", u.listAllUsers)
		r.Get("/{id}", u.getUser)
		r.Put("/{id}", u.updateUser)
		r.Delete("/{id}", u.deleteUser)
	})
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var mu models.User

	if err := utils.DecodeJSON(r, &mu); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid json payload")
		return
	}

	mu.ID = ""

	userCreated, err := u.db.Insert(models.User{
		FirstName: mu.FirstName,
		LastName:  mu.LastName,
		Biography: mu.Biography,
	})
	if err != nil {
		if errors.Is(err, store.ErrBadInput) {
			utils.WriteError(w, http.StatusBadRequest, "first_name, last_name and biography are required and must respect size limits")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save user")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "User created successfully", userCreated)
}

func (u *UserHandler) listAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.db.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Users fetched successfully", users)
}

func (u *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := u.db.FindById(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User fetched successfully", user)
}

func (u *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var mu models.User
	if err := utils.DecodeJSON(r, &mu); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid json payload")
		return
	}

	updateUser, err := u.db.Update(id, models.User{
		FirstName: mu.FirstName,
		LastName:  mu.LastName,
		Biography: mu.Biography,
	})
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			utils.WriteError(w, http.StatusNotFound, "User not found")
		case errors.Is(err, store.ErrBadInput):
			utils.WriteError(w, http.StatusBadRequest, "first_name, last_name and biography are required and must respect size limits")
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User updated successfully", updateUser)
}

func (u *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := u.db.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User deleted successfully", user)
}
