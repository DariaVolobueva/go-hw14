package service

import (
	"api/model"
	"api/repository"
	"log"
	"time"
)

type TourService struct {
	repo *repository.TourRepository
}

func NewTourService(repo *repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) GetAvailableTours() []model.Tour {
	return s.repo.GetAllTours()
}

func (s *TourService) BookTour(tourID int, email string) error {
	booking := model.Booking{
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

func (s *TourService) getTourByID(id int) model.Tour {
	for _, tour := range s.repo.GetAllTours() {
		if tour.ID == id {
			return tour
		}
	}
	return model.Tour{}
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