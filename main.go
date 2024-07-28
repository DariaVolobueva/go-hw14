package main

import (
	"api/handler"
	"api/repository"
	"api/service"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewTourRepository()
	service := service.NewTourService(repo)
	handler := handler.NewTourHandler(service)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/tours", "/book", "/bookings":
			handler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}