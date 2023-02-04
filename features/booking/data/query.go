package data

import (
	"campyuk-api/features/booking"

	"gorm.io/gorm"
)

type bookingData struct {
	db *gorm.DB
}

func New(db *gorm.DB) booking.BookingData {
	return &bookingData{
		db: db,
	}
}

func (bd *bookingData) Create(userID uint, virtualNumber string) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bd *bookingData) Update(bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bd *bookingData) List(userID uint) ([]booking.Core, error) {
	return []booking.Core{}, nil
}

func (bd *bookingData) GetByID(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bd *bookingData) Cancel(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}
