package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Handle /users
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"message": "Hello from Users Service",
			"users": ["Alice", "Bob", "Charlie"]
		}`))
	})

	// Also support root for health
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Users service is running"))
	})

	http.ListenAndServe(":8081", r)
}