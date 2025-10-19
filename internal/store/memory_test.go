package store_test

import (
	"testing"

	"aantonioprado/rs-go-api-crud-memory/internal/models"
	"aantonioprado/rs-go-api-crud-memory/internal/store"
)

func TestMemory_CRUD(t *testing.T) {
	mem := store.NewMemory()

	u, err := mem.Insert(models.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Biography: "Biografia válida com pelo menos vinte caracteres.",
	})
	if err != nil {
		t.Fatalf("insert: unexpected err: %v", err)
	}
	if u.ID == "" {
		t.Fatalf("insert: expected UUID to be generated")
	}

	got, err := mem.FindById(u.ID)
	if err != nil {
		t.Fatalf("findById: unexpected err: %v", err)
	}
	if got.FirstName != "Jane" {
		t.Fatalf("findById: wrong first name: %s", got.FirstName)
	}

	u2, err := mem.Update(u.ID, models.User{
		FirstName: "Janette",
		LastName:  "Doe",
		Biography: "Biografia atualizada com pelo menos vinte caracteres.",
	})
	if err != nil {
		t.Fatalf("update: unexpected err: %v", err)
	}
	if u2.FirstName != "Janette" {
		t.Fatalf("update: did not update FirstName")
	}

	del, err := mem.Delete(u.ID)
	if err != nil {
		t.Fatalf("delete: unexpected err: %v", err)
	}
	if del.ID != u.ID {
		t.Fatalf("delete: returned different user")
	}

	if _, err := mem.FindById(u.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}

func TestMemory_Validation(t *testing.T) {
	mem := store.NewMemory()
	_, err := mem.Insert(models.User{
		FirstName: "A",
		LastName:  "Ok",
		Biography: "Biografia válida com pelo menos vinte caracteres.",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
