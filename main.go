package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Tour struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Transport   string  `json:"transport"`
}

type Booking struct {
	ID     int       `json:"id"`
	TourID int       `json:"tour_id"`
	Email  string    `json:"email"`
	Date   time.Time `json:"date"`
}

type TourRepository struct {
	tours    []Tour
	bookings []Booking
}

func NewTourRepository() *TourRepository {
	return &TourRepository{
		tours: []Tour{
			{ID: 1, Name: "Paris Adventure", Description: "Explore the City of Light", Price: 1000, Transport: "Plane"},
			{ID: 2, Name: "Rome Getaway", Description: "Discover ancient history", Price: 1200, Transport: "Train"},
		},
		bookings: []Booking{},
	}
}

func (r *TourRepository) GetAllTours() []Tour {
	return r.tours
}

func (r *TourRepository) AddBooking(booking Booking) {
	booking.ID = len(r.bookings) + 1
	r.bookings = append(r.bookings, booking)
}

func (r *TourRepository) GetAllBookings() []Booking {
	return r.bookings
}

type TourService struct {
	repo *TourRepository
}

func NewTourService(repo *TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) GetAvailableTours() []Tour {
	return s.repo.GetAllTours()
}

func (s *TourService) BookTour(tourID int, email string) error {
	booking := Booking{
		TourID: tourID,
		Email:  email,
		Date:   time.Now(),
	}
	s.repo.AddBooking(booking)
	s.sendConfirmationEmail(email)
	return nil
}

func (s *TourService) GetBookedTours() []map[string]interface{} {
	bookings := s.repo.GetAllBookings()
	result := make([]map[string]interface{}, len(bookings))

	for i, booking := range bookings {
		tour := s.getTourByID(booking.TourID)
		status := s.calculateTourStatus(booking.Date)

		result[i] = map[string]interface{}{
			"id":     booking.ID,
			"tour":   tour,
			"email":  booking.Email,
			"date":   booking.Date,
			"status": status,
		}
	}

	return result
}

func (s *TourService) getTourByID(id int) Tour {
	for _, tour := range s.repo.GetAllTours() {
		if tour.ID == id {
			return tour
		}
	}
	return Tour{}
}

func (s *TourService) calculateTourStatus(bookingDate time.Time) string {
	now := time.Now()
	diff := bookingDate.Sub(now)

	switch {
	case diff > 30*24*time.Hour:
		return "Майбутній"
	case diff > 7*24*time.Hour:
		return "Скоро"
	case diff > 0:
		return "Триває"
	default:
		return "Завершився"
	}
}

func (s *TourService) sendConfirmationEmail(email string) {
	log.Printf("Відправлено підтвердження на email: %s", email)
}

type TourHandler struct {
	service *TourService
}

func NewTourHandler(service *TourService) *TourHandler {
	return &TourHandler{service: service}
}

func (h *TourHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/tours":
		h.GetAvailableTours(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/book":
		h.BookTour(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/bookings":
		h.GetBookedTours(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *TourHandler) GetAvailableTours(w http.ResponseWriter, r *http.Request) {
	tours := h.service.GetAvailableTours()
	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) BookTour(w http.ResponseWriter, r *http.Request) {
	var booking struct {
		TourID int    `json:"tour_id"`
		Email  string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.BookTour(booking.TourID, booking.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Тур успішно заброньовано")
}

func (h *TourHandler) GetBookedTours(w http.ResponseWriter, r *http.Request) {
	bookedTours := h.service.GetBookedTours()
	json.NewEncoder(w).Encode(bookedTours)
}

func main() {
	repo := NewTourRepository()
	service := NewTourService(repo)
	handler := NewTourHandler(service)

	http.Handle("/tours", handler)
	http.Handle("/book", handler)
	http.Handle("/bookings", handler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}