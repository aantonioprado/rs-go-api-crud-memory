package store

import (
	"aantonioprado/rs-go-api-crud-memory/internal/models"
	"errors"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("User not found")
	ErrBadInput = errors.New("Bad Request")
)

type Memory struct {
	mux  sync.RWMutex
	data map[string]models.User
}

func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]models.User),
	}
}

func validate(u models.User) error {
	fn := strings.TrimSpace(u.FirstName)
	ln := strings.TrimSpace(u.LastName)
	bio := strings.TrimSpace(u.Biography)

	if fn == "" || ln == "" || bio == "" {
		return ErrBadInput
	}
	if n := utf8.RuneCountInString(fn); n < 2 || n > 20 {
		return ErrBadInput
	}
	if n := utf8.RuneCountInString(ln); n < 2 || n > 20 {
		return ErrBadInput
	}
	if n := utf8.RuneCountInString(bio); n < 20 || n > 450 {
		return ErrBadInput
	}
	return nil
}

func (m *Memory) FindAll() ([]models.User, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	out := make([]models.User, 0, len(m.data))
	for _, u := range m.data {
		out = append(out, u)
	}
	return out, nil
}

func (m *Memory) FindById(id string) (*models.User, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	u, ok := m.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &u, nil
}

func (m *Memory) Insert(newUser models.User) (models.User, error) {
	if err := validate(newUser); err != nil {
		return models.User{}, err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	newUser.ID = uuid.NewString()
	m.data[newUser.ID] = newUser

	return newUser, nil
}

func (m *Memory) Update(id string, updateUser models.User) (models.User, error) {
	if err := validate(updateUser); err != nil {
		return models.User{}, err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	if _, ok := m.data[id]; !ok {
		return models.User{}, ErrNotFound
	}

	updateUser.ID = id
	m.data[id] = updateUser

	return updateUser, nil
}

func (m *Memory) Delete(id string) (models.User, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	u, ok := m.data[id]
	if !ok {
		return models.User{}, ErrNotFound
	}

	delete(m.data, id)
	return u, nil
}
