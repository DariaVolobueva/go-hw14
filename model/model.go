package model

import "time"

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