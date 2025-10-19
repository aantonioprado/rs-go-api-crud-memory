package main

import (
	"aantonioprado/rs-go-api-crud-memory/internal/api"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3200"
	}

	router := api.NewRouter()

	log.Printf("API running on http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
