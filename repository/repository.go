package repository

import (
	"api/model"
)

type TourRepository struct {
	tours    []model.Tour
	bookings []model.Booking
}

func NewTourRepository() *TourRepository {
	return &TourRepository{
		tours: []model.Tour{
			{ID: 1, Name: "Paris Adventure", Description: "Explore the City of Light", Price: 1000, Transport: "Plane"},
			{ID: 2, Name: "Rome Getaway", Description: "Discover ancient history", Price: 1200, Transport: "Train"},
		},
		bookings: []model.Booking{},
	}
}

func (r *TourRepository) GetAllTours() []model.Tour {
	return r.tours
}

func (r *TourRepository) AddBooking(booking model.Booking) {
	booking.ID = len(r.bookings) + 1
	r.bookings = append(r.bookings, booking)
}

func (r *TourRepository) GetAllBookings() []model.Booking {
	return r.bookings
}