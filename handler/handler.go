package handler

import (
	"api/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type TourHandler struct {
	service *service.TourService
}

func NewTourHandler(service *service.TourService) *TourHandler {
	return &TourHandler{service: service}
}

func (h *TourHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/tours":
		h.GetAvailableTours(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/book":
		h.BookTour(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/book":
		h.BookingInfo(w, r)
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

func (h *TourHandler) BookingInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Для бронювання туру використовуйте POST запит з параметрами tour_id та email",
		"example": `{"tour_id": 1, "email": "user@example.com"}`,
	})
}

func (h *TourHandler) GetBookedTours(w http.ResponseWriter, r *http.Request) {
	bookedTours := h.service.GetBookedTours()
	json.NewEncoder(w).Encode(bookedTours)
}